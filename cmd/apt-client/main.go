package main

import (
    "flag"
    "io"
    "log"
    "net"

    "github.com/meislaozhang/apt-proxy/internal/client"
    "github.com/meislaozhang/apt-proxy/internal/proxy"
)

func main() {
    listen := flag.String("listen", "127.0.0.1:1080", "SOCKS5 listen address")
    remote := flag.String("server", "127.0.0.1:8443", "APT server address")
    name := flag.String("server-name", "localhost", "TLS ServerName")
    token := flag.String("token", "change-me", "authentication token")
    insecure := flag.Bool("insecure-skip-verify", true, "skip certificate verification; development only")
    flag.Parse()

    s, err := client.Dial(client.Config{ServerAddr: *remote, ServerName: *name, Token: *token, InsecureSkipVerify: *insecure})
    if err != nil { log.Fatal(err) }
    defer s.Close()

    ln, err := net.Listen("tcp", *listen)
    if err != nil { log.Fatal(err) }
    defer ln.Close()
    log.Printf("APT SOCKS5 listening on %s (shared multiplexed session)", *listen)

    err = proxy.ServeSOCKS5(ln, func(c net.Conn, target string) error {
        r, err := s.Open(target)
        if err != nil { return err }
        defer r.Close()
        go func() { _, _ = io.Copy(r, c); _ = r.Close() }()
        _, _ = io.Copy(c, r)
        return nil
    })
    log.Fatal(err)
}
