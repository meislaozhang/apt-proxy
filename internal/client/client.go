package client

import (
    "bufio"
    "crypto/tls"
    "fmt"
    "io"
    "net"
    "sync"
    "time"

    "github.com/meislaozhang/apt-proxy/internal/protocol"
)

type Config struct { ServerAddr, ServerName, Token string; InsecureSkipVerify bool }

type Session struct { c *tls.Conn; r *bufio.Reader; mu sync.Mutex }

func Dial(cfg Config) (*Session, error) {
    c, err := tls.DialWithDialer(&net.Dialer{Timeout:10*time.Second}, "tcp", cfg.ServerAddr, &tls.Config{MinVersion:tls.VersionTLS13, ServerName:cfg.ServerName, InsecureSkipVerify:cfg.InsecureSkipVerify})
    if err != nil { return nil, err }
    s := &Session{c:c, r:bufio.NewReader(c)}
    if err := s.write(protocol.Frame{Type:protocol.TypeAuth, Payload:[]byte(cfg.Token)}); err != nil { _=c.Close(); return nil,err }
    f,err := protocol.ReadFrame(s.r); if err != nil || f.Type != protocol.TypeAuthOK { _=c.Close(); return nil,fmt.Errorf("authentication failed") }
    return s,nil
}
func (s *Session) write(f protocol.Frame) error { s.mu.Lock(); defer s.mu.Unlock(); return f.Write(s.c) }
func (s *Session) Close() error { return s.c.Close() }

func (s *Session) Open(target string) (net.Conn,error) {
    if err:=s.write(protocol.Frame{Type:protocol.TypeOpen,StreamID:1,Payload:[]byte(target)}); err!=nil { return nil,err }
    for { f,err:=protocol.ReadFrame(s.r); if err!=nil{return nil,err}; if f.StreamID!=1{continue}; if f.Type==protocol.TypeOpenOK{return &stream{sess:s},nil}; return nil,fmt.Errorf("open failed: %s",f.Payload) }
}

type stream struct { sess *Session; mu sync.Mutex; closed bool }
func (s *stream) Read(p []byte)(int,error){ for { f,e:=protocol.ReadFrame(s.sess.r); if e!=nil{return 0,e}; if f.StreamID!=1{continue}; switch f.Type{case protocol.TypeData:return copy(p,f.Payload),nil;case protocol.TypeClose:return 0,io.EOF;case protocol.TypeError:return 0,fmt.Errorf("remote: %s",f.Payload)} } }
func (s *stream) Write(p []byte)(int,error){ if err:=s.sess.write(protocol.Frame{Type:protocol.TypeData,StreamID:1,Payload:append([]byte(nil),p...)});err!=nil{return 0,err};return len(p),nil }
func (s *stream) Close()error{s.mu.Lock();defer s.mu.Unlock();if s.closed{return nil};s.closed=true;return s.sess.write(protocol.Frame{Type:protocol.TypeClose,StreamID:1})}
func(s *stream)LocalAddr()net.Addr{return s.sess.c.LocalAddr()};func(s *stream)RemoteAddr()net.Addr{return s.sess.c.RemoteAddr()};func(s *stream)SetDeadline(t time.Time)error{return s.sess.c.SetDeadline(t)};func(s *stream)SetReadDeadline(t time.Time)error{return s.sess.c.SetReadDeadline(t)};func(s *stream)SetWriteDeadline(t time.Time)error{return s.sess.c.SetWriteDeadline(t)}
