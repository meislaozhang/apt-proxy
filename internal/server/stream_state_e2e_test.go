package server

import (
    "bufio"
    "crypto/tls"
    "fmt"
    "io"
    "net"
    "os"
    "testing"
    "time"

    "github.com/meislaozhang/apt-proxy/internal/client"
)

func startE2EServer(t *testing.T, token string) (*client.Session, func()) {
    t.Helper()
    certFile, keyFile := testCertificate(t)
    ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{MinVersion: tls.VersionTLS13, Certificates: loadTestCertificate(t, certFile, keyFile)})
    if err != nil { t.Fatal(err) }
    srv := New(Config{Addr: ln.Addr().String(), Token: token})
    stop := make(chan struct{})
    go func() {
        for {
            c, err := ln.Accept()
            if err != nil { return }
            go srv.handle(c)
        }
    }()
    sess, err := client.Dial(client.Config{ServerAddr: ln.Addr().String(), ServerName: "localhost", Token: token, InsecureSkipVerify: true})
    if err != nil { ln.Close(); t.Fatal(err) }
    cleanup := func() { sess.Close(); close(stop); ln.Close(); os.Remove(certFile); os.Remove(keyFile) }
    return sess, cleanup
}

func TestE2EHalfClose(t *testing.T) {
    target, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil { t.Fatal(err) }
    defer target.Close()
    go func() {
        c, err := target.Accept()
        if err != nil { return }
        defer c.Close()
        r := bufio.NewReader(c)
        _, _ = r.ReadString('\n')
        // The target only responds after observing EOF on the request direction.
        b, _ := io.ReadAll(r)
        _ = b
        _, _ = c.Write([]byte("half-close-ok"))
    }()

    sess, cleanup := startE2EServer(t, "half-close-token")
    defer cleanup()
    conn, err := sess.Open(target.Addr().String())
    if err != nil { t.Fatal(err) }
    if _, err := conn.Write([]byte("hello\n")); err != nil { t.Fatal(err) }
    cw, ok := conn.(interface{ CloseWrite() error })
    if !ok { t.Fatal("stream does not expose CloseWrite") }
    if err := cw.CloseWrite(); err != nil { t.Fatal(err) }
    _ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
    got, err := io.ReadAll(conn)
    if err != nil { t.Fatal(err) }
    if string(got) != "half-close-ok" { t.Fatalf("got %q", got) }
}

func TestE2EResetKeepsSessionAlive(t *testing.T) {
    target, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil { t.Fatal(err) }
    defer target.Close()
    go func() {
        for {
            c, err := target.Accept()
            if err != nil { return }
            go func(c net.Conn) { defer c.Close(); _, _ = io.Copy(c, c) }(c)
        }
    }()

    sess, cleanup := startE2EServer(t, "reset-token")
    defer cleanup()
    first, err := sess.Open(target.Addr().String())
    if err != nil { t.Fatal(err) }
    r, ok := first.(interface{ Reset() error })
    if !ok { t.Fatal("stream does not expose Reset") }
    if err := r.Reset(); err != nil { t.Fatal(err) }

    second, err := sess.Open(target.Addr().String())
    if err != nil { t.Fatal(fmt.Errorf("session unusable after stream reset: %w", err)) }
    defer second.Close()
    if _, err := second.Write([]byte("still-alive")); err != nil { t.Fatal(err) }
    _ = second.SetReadDeadline(time.Now().Add(3 * time.Second))
    got := make([]byte, len("still-alive"))
    if _, err := io.ReadFull(second, got); err != nil { t.Fatal(err) }
    if string(got) != "still-alive" { t.Fatalf("got %q", got) }
}
