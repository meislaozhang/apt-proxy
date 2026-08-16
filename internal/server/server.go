package server

import (
    "bufio"
    "crypto/tls"
    "net"
    "time"

    "github.com/meislaozhang/apt-proxy/internal/auth"
    "github.com/meislaozhang/apt-proxy/internal/protocol"
)

type Config struct { Addr, CertFile, KeyFile, Token string }
type Server struct { cfg Config }

func New(cfg Config) *Server { return &Server{cfg: cfg} }

func (s *Server) ListenAndServe() error {
    cert, err := tls.LoadX509KeyPair(s.cfg.CertFile, s.cfg.KeyFile); if err != nil { return err }
    ln, err := tls.Listen("tcp", s.cfg.Addr, &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS13}); if err != nil { return err }
    defer ln.Close()
    for { c, err := ln.Accept(); if err != nil { return err }; go s.handle(c) }
}

func (s *Server) handle(c net.Conn) {
    defer c.Close()
    _ = c.SetDeadline(time.Now().Add(15*time.Second))
    r := bufio.NewReader(c)
    f, err := protocol.ReadFrame(r); if err != nil { return }
    if f.Type != protocol.TypeAuth || !auth.EqualToken(f.Payload, []byte(s.cfg.Token)) { return }
    if err = (protocol.Frame{Type: protocol.TypeAuthOK}).Write(c); err != nil { return }
    f, err = protocol.ReadFrame(r); if err != nil || f.Type != protocol.TypeOpen { return }
    _ = c.SetDeadline(time.Time{})
    id := f.StreamID
    d, err := net.DialTimeout("tcp", string(f.Payload), 10*time.Second)
    if err != nil { _ = (protocol.Frame{Type: protocol.TypeError, StreamID:id, Payload:[]byte(err.Error())}).Write(c); return }
    defer d.Close()
    if err = (protocol.Frame{Type:protocol.TypeOpenOK, StreamID:id}).Write(c); err != nil { return }
    done := make(chan struct{})
    go func() {
        defer close(done)
        buf := make([]byte, 32*1024)
        for { n,e := d.Read(buf); if n>0 { if we := (protocol.Frame{Type:protocol.TypeData, StreamID:id, Payload:append([]byte(nil),buf[:n]...)}).Write(c); we != nil { return } }; if e != nil { _=(protocol.Frame{Type:protocol.TypeClose,StreamID:id}).Write(c); return } }
    }()
    for { f,err = protocol.ReadFrame(r); if err != nil { return }; if f.StreamID != id { continue }; switch f.Type { case protocol.TypeData: if _,err=d.Write(f.Payload); err != nil { return }; case protocol.TypeClose: return } }
    <-done
}
