package server

import (
    "bufio"
    "context"
    "crypto/tls"
    "net"
    "sync"
    "time"
    "github.com/meislaozhang/apt-proxy/internal/auth"
    "github.com/meislaozhang/apt-proxy/internal/protocol"
)

type Config struct { Addr, CertFile, KeyFile, Token string }
type Server struct { cfg Config }
type serverSession struct { c net.Conn; r *bufio.Reader; writeMu sync.Mutex; mu sync.Mutex; streams map[uint32]net.Conn; send *protocol.BlockingFlowWindow }
const ( defaultSessionWindow uint64 = 64*1024; defaultStreamWindow uint64 = 64*1024; maxDataChunk uint64 = 16*1024 )
func New(cfg Config)*Server{return &Server{cfg:cfg}}
func(s *Server)ListenAndServe()error{cert,err:=tls.LoadX509KeyPair(s.cfg.CertFile,s.cfg.KeyFile);if err!=nil{return err};ln,err:=tls.Listen("tcp",s.cfg.Addr,&tls.Config{Certificates:[]tls.Certificate{cert},MinVersion:tls.VersionTLS13});if err!=nil{return err};defer ln.Close();for{c,err:=ln.Accept();if err!=nil{return err};go s.handle(c)}}
func(s *Server)handle(c net.Conn){send,err:=protocol.NewBlockingFlowWindow(defaultSessionWindow,defaultSessionWindow);if err!=nil{_=c.Close();return};ss:=&serverSession{c:c,r:bufio.NewReader(c),streams:make(map[uint32]net.Conn),send:send};defer ss.closeAll();_=c.SetDeadline(time.Now().Add(15*time.Second));f,err:=protocol.ReadFrame(ss.r);if err!=nil{return};if f.Type!=protocol.TypeAuth||!auth.EqualToken(f.Payload,[]byte(s.cfg.Token)){return};if err=(protocol.Frame{Type:protocol.TypeAuthOK}).Write(c);err!=nil{return};_=c.SetDeadline(time.Time{});for{f,err=protocol.ReadFrame(ss.r);if err!=nil{return};switch f.Type{case protocol.TypeOpen:ss.open(f.StreamID,string(f.Payload));case protocol.TypeData:ss.writeData(f.StreamID,f.Payload);case protocol.TypeWindowUpdate:ss.windowUpdate(f.StreamID,f.Payload);case protocol.TypeHalfClose:ss.halfClose(f.StreamID);case protocol.TypeReset:ss.reset(f.StreamID);case protocol.TypeClose:ss.closeStream(f.StreamID)}}}
func(ss *serverSession)write(f protocol.Frame)error{ss.writeMu.Lock();defer ss.writeMu.Unlock();return f.Write(ss.c)}
func(ss *serverSession)windowUpdate(id uint32,p []byte){n,ok:=protocol.DecodeWindowUpdate(p);if !ok{return};if id==0{_=ss.send.Add(n)}}
func(ss *serverSession)open(id uint32,target string){ss.mu.Lock();if _,ok:=ss.streams[id];ok{ss.mu.Unlock();return};ss.mu.Unlock();d,err:=net.DialTimeout("tcp",target,10*time.Second);if err!=nil{_=ss.write(protocol.Frame{Type:protocol.TypeError,StreamID:id,Payload:[]byte(err.Error())});return};ss.mu.Lock();ss.streams[id]=d;ss.mu.Unlock();if err:=ss.write(protocol.Frame{Type:protocol.TypeOpenOK,StreamID:id});err!=nil{ss.closeStream(id);return};go ss.pump(id,d)}
func(ss *serverSession)pump(id uint32,d net.Conn){buf:=make([]byte,32*1024);for{n,err:=d.Read(buf);if n>0{off:=0;for off<n{chunk:=uint64(n-off);if chunk>maxDataChunk{chunk=maxDataChunk};if e:=ss.send.Acquire(context.Background(),chunk);e!=nil{ss.closeStream(id);return};if we:=ss.write(protocol.Frame{Type:protocol.TypeData,StreamID:id,Payload:append([]byte(nil),buf[off:off+int(chunk)]...)});we!=nil{ss.closeStream(id);return};off+=int(chunk)}};if err!=nil{ss.closeStream(id);return}}}
func(ss *serverSession)writeData(id uint32,p []byte){ss.mu.Lock();d:=ss.streams[id];ss.mu.Unlock();if d==nil{return};if _,err:=d.Write(p);err!=nil{ss.closeStream(id);return};_=ss.write(protocol.Frame{Type:protocol.TypeWindowUpdate,StreamID:id,Payload:protocol.EncodeWindowUpdate(uint64(len(p)))});_=ss.write(protocol.Frame{Type:protocol.TypeWindowUpdate,StreamID:0,Payload:protocol.EncodeWindowUpdate(uint64(len(p)))})}
func(ss *serverSession)halfClose(id uint32){ss.mu.Lock();d:=ss.streams[id];ss.mu.Unlock();if d==nil{return};if cw,ok:=d.(interface{CloseWrite()error});ok{_=cw.CloseWrite()}}
func(ss *serverSession)reset(id uint32){ss.mu.Lock();d:=ss.streams[id];delete(ss.streams,id);ss.mu.Unlock();if d!=nil{_=d.Close()};_=ss.write(protocol.Frame{Type:protocol.TypeReset,StreamID:id})}
func(ss *serverSession)closeStream(id uint32){ss.mu.Lock();d:=ss.streams[id];delete(ss.streams,id);ss.mu.Unlock();if d!=nil{_=d.Close()};_=ss.write(protocol.Frame{Type:protocol.TypeClose,StreamID:id})}
func(ss *serverSession)closeAll(){ss.mu.Lock();for id,d:=range ss.streams{_=d.Close();delete(ss.streams,id)};ss.mu.Unlock();_=ss.c.Close()}
