package main

import (
    "flag"
    "io"
    "log"
    "net"
    "github.com/meislaozhang/apt-proxy/internal/client"
    "github.com/meislaozhang/apt-proxy/internal/proxy"
)

func main(){
    listen:=flag.String("listen","127.0.0.1:1080","SOCKS5 listen")
    remote:=flag.String("server","127.0.0.1:8443","APT server")
    name:=flag.String("server-name","localhost","TLS server name")
    token:=flag.String("token","change-me","authentication token")
    insecure:=flag.Bool("insecure-skip-verify",true,"skip certificate verification for development only")
    flag.Parse()
    ln,err:=net.Listen("tcp",*listen);if err!=nil{log.Fatal(err)}
    log.Printf("APT SOCKS5 listening on %s",*listen)
    err=proxy.ServeSOCKS5(ln,func(c net.Conn,target string)error{
        s,e:=client.Dial(client.Config{ServerAddr:*remote,ServerName:*name,Token:*token,InsecureSkipVerify:*insecure});if e!=nil{return e};defer s.Close()
        r,e:=s.Open(target);if e!=nil{return e};defer r.Close()
        go func(){_,_=io.Copy(r,c);_ = r.Close()}()
        _,_=io.Copy(c,r);return nil
    })
    log.Fatal(err)
}
