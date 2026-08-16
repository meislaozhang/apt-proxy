package main

import (
    "flag"
    "log"
    "github.com/meislaozhang/apt-proxy/internal/server"
)

func main() {
    addr:=flag.String("listen",":8443","listen address")
    cert:=flag.String("cert","server.crt","TLS certificate")
    key:=flag.String("key","server.key","TLS private key")
    token:=flag.String("token","change-me","authentication token")
    flag.Parse()
    log.Printf("APT server listening on %s",*addr)
    log.Fatal(server.New(server.Config{Addr:*addr,CertFile:*cert,KeyFile:*key,Token:*token}).ListenAndServe())
}
