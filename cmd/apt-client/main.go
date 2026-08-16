package main

import (
    "io"
    "log"
    "net"

    "github.com/meislaozhang/apt-proxy/internal/client"
    "github.com/meislaozhang/apt-proxy/internal/proxy"
)

func main() {
    listen := flagString("listen", "127.0.0.1:1080", "SOCKS5 listen address")
    remote := flagString("server", "127.0.0.1:8443", "APT server address")
    name := flagString("server-name", "localhost", "TLS ServerName")
    token := flagString("token", "change-me", "authentication token")
    insecure := flagBool("insecure-skip-verify", true, "skip certificate verification; development only")
    flagParse()

    s, err := client.Dial(client.Config{
        ServerAddr: remoteValue(remote), ServerName: remoteValue(name),
        Token: remoteValue(token), InsecureSkipVerify: boolValue(insecure),
    })
    if err != nil { log.Fatal(err) }
    defer s.Close()

    ln, err := net.Listen("tcp", remoteValue(listen))
    if err != nil { log.Fatal(err) }
    defer ln.Close()
    log.Printf("APT SOCKS5 listening on %s (shared multiplexed session)", remoteValue(listen))

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

// Small wrappers keep the command implementation explicit while avoiding a global flag package dependency in tests.
// They are backed by the standard flag package at runtime.

type stringFlag struct { value string }
type boolFlag struct { value bool }

func flagString(name, value, usage string) *stringFlag { return &stringFlag{value: registerString(name, value, usage)} }
func flagBool(name string, value bool, usage string) *boolFlag { return &boolFlag{value: registerBool(name, value, usage)} }
func flagParse() { parseFlags() }
func remoteValue(v *stringFlag) string { return v.value }
func boolValue(v *boolFlag) bool { return v.value }
