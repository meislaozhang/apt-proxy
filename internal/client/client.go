package client

import (
    "bufio"
    "context"
    "crypto/tls"
    "encoding/binary"
    "fmt"
    "io"
    "net"
    "sync"
    "sync/atomic"
    "time"

    "github.com/meislaozhang/apt-proxy/internal/protocol"
)

const (
    defaultSessionWindow uint64 = 64 * 1024
    defaultStreamWindow  uint64 = 64 * 1024
    maxDataChunk         uint64 = 16 * 1024
)

type Config struct { ServerAddr, ServerName, Token string; InsecureSkipVerify bool }
type Session struct { c *tls.Conn; r *bufio.Reader; writeMu sync.Mutex; mu sync.Mutex; streams map[uint32]*stream; nextID uint32; done chan struct{}; once sync.Once; send *sendAdmission }
type stream struct { sess *Session; id uint32; data chan []byte; openOK chan error; closeOnce sync.Once; closed chan struct{}; readMu sync.Mutex; pending []byte; localWriteClosed atomic.Bool; send *protocol.BlockingFlowWindow }

func Dial(cfg Config) (*Session,error) {
    c,err:=tls.DialWithDialer(&net.Dialer{Timeout:10*time.Second},"tcp",cfg.ServerAddr,&tls.Config{MinVersion:tls.VersionTLS13,ServerName:cfg.ServerName,InsecureSkipVerify:cfg.InsecureSkipVerify})
    if err!=nil{return nil,err}
    send, err := newSendAdmission(defaultSessionWindow, defaultStreamWindow)
    if err != nil { _ = c.Close(); return nil, err }
    s:=&Session{c:c,r:bufio.NewReader(c),streams:make(map[uint32]*stream),nextID:1,done:make(chan struct{}),send:send}
    if err:=s.write(protocol.Frame{Type:protocol.TypeAuth,Payload:[]byte(cfg.Token)});err!=nil{_ = c.Close();return nil,err}
    f,err:=protocol.ReadFrame(s.r);if err!=nil||f.Type!=protocol.TypeAuthOK{_ = c.Close();return nil,fmt.Errorf("authentication failed")}
    go s.readLoop();return s,nil
}

func(s *Session)write(f protocol.Frame)error{s.writeMu.Lock();defer s.writeMu.Unlock();return f.Write(s.c)}

func(s *Session)readLoop(){
    defer s.shutdown()
    for{
        f,err:=protocol.ReadFrame(s.r);if err!=nil{return}
        if f.Type == protocol.TypeWindowUpdate {
            n, ok := decodeWindowUpdate(f.Payload)
            if !ok || n == 0 { continue }
            if f.StreamID == 0 { _ = s.send.updateSession(n) } else {
                s.mu.Lock(); st:=s.streams[f.StreamID]; s.mu.Unlock()
                if st != nil { _ = st.send.Add(n) }
            }
            continue
        }
        s.mu.Lock();st:=s.streams[f.StreamID];s.mu.Unlock();if st==nil{continue}
        switch f.Type{case protocol.TypeOpenOK:select{case st.openOK<-nil:default:};case protocol.TypeError:select{case st.openOK<-fmt.Errorf("remote: %s",f.Payload):default:};case protocol.TypeData:b:=append([]byte(nil),f.Payload...);select{case st.data<-b:case <-st.closed:};case protocol.TypeClose,protocol.TypeReset:st.closeOnce.Do(func(){close(st.closed)})}
    }
}

func decodeWindowUpdate(p []byte) (uint64, bool) {
    if len(p) != 8 { return 0, false }
    return binary.BigEndian.Uint64(p), true
}

func(s *Session)shutdown(){s.once.Do(func(){close(s.done);_=s.c.Close();s.mu.Lock();defer s.mu.Unlock();for _,st:=range s.streams{st.closeOnce.Do(func(){close(st.closed)})}})}
func(s *Session)Close()error{s.shutdown();return nil}

func(s *Session)Open(target string)(net.Conn,error){
    id:=atomic.AddUint32(&s.nextID,1)
    streamWindow, err := protocol.NewBlockingFlowWindow(defaultStreamWindow, defaultStreamWindow)
    if err != nil { return nil, err }
    st:=&stream{sess:s,id:id,data:make(chan []byte,16),openOK:make(chan error,1),closed:make(chan struct{}),send:streamWindow}
    s.mu.Lock();s.streams[id]=st;s.mu.Unlock()
    if err:=s.write(protocol.Frame{Type:protocol.TypeOpen,StreamID:id,Payload:[]byte(target)});err!=nil{s.remove(id);return nil,err}
    select{case err:=<-st.openOK:if err!=nil{s.remove(id);return nil,err};return st,nil;case <-s.done:s.remove(id);return nil,io.EOF;case <-time.After(15*time.Second):s.remove(id);return nil,fmt.Errorf("open timeout")}
}
func(s *Session)remove(id uint32){s.mu.Lock();delete(s.streams,id);s.mu.Unlock()}

func(st *stream)Read(p []byte)(int,error){st.readMu.Lock();defer st.readMu.Unlock();for len(st.pending)==0{select{case b:=<-st.data:st.pending=b;continue;default:};select{case b:=<-st.data:st.pending=b;case <-st.closed:return 0,io.EOF;case <-st.sess.done:return 0,io.EOF}};n:=copy(p,st.pending);st.pending=st.pending[n:];return n,nil}

func(st *stream)Write(p []byte)(int,error){
    if st.localWriteClosed.Load(){return 0,io.ErrClosedPipe}
    select{case <-st.closed:return 0,io.ErrClosedPipe;default:}
    written := 0
    for written < len(p) {
        n := len(p) - written
        if uint64(n) > maxDataChunk { n = int(maxDataChunk) }
        ctx, cancel := context.WithCancel(context.Background())
        go func(){ select { case <-st.closed: cancel(); case <-st.sess.done: cancel(); case <-ctx.Done(): } }()
        if err:=st.send.Acquire(ctx,uint64(n));err!=nil{cancel();return written,err}
        if err:=st.sess.send.acquire(ctx,uint64(n));err!=nil{_ = st.send.Add(uint64(n));cancel();return written,err}
        cancel()
        if err:=st.sess.write(protocol.Frame{Type:protocol.TypeData,StreamID:st.id,Payload:append([]byte(nil),p[written:written+n]...)});err!=nil{return written,err}
        written += n
    }
    return written,nil
}

func(st *stream)CloseWrite()error{if st.localWriteClosed.Swap(true){return nil};return st.sess.write(protocol.Frame{Type:protocol.TypeHalfClose,StreamID:st.id})}
func(st *stream)Reset()error{st.closeOnce.Do(func(){close(st.closed);_=st.sess.write(protocol.Frame{Type:protocol.TypeReset,StreamID:st.id});st.sess.remove(st.id)});return nil}
func(st *stream)Close()error{st.closeOnce.Do(func(){close(st.closed);_=st.sess.write(protocol.Frame{Type:protocol.TypeClose,StreamID:st.id});st.sess.remove(st.id)});return nil}
func(st *stream)LocalAddr()net.Addr{return st.sess.c.LocalAddr()};func(st *stream)RemoteAddr()net.Addr{return st.sess.c.RemoteAddr()};func(st *stream)SetDeadline(t time.Time)error{return st.sess.c.SetDeadline(t)};func(st *stream)SetReadDeadline(t time.Time)error{return st.sess.c.SetReadDeadline(t)};func(st *stream)SetWriteDeadline(t time.Time)error{return st.sess.c.SetWriteDeadline(t)}
