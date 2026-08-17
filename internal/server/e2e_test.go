package server

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/tls"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "fmt"
    "io"
    "math/big"
    "net"
    "os"
    "sync"
    "testing"
    "time"

    "github.com/meislaozhang/apt-proxy/internal/client"
)

func TestE2ETCPMultiplexing(t *testing.T) {
    certFile, keyFile := testCertificate(t)
    defer os.Remove(certFile)
    defer os.Remove(keyFile)

    echoLn, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil { t.Fatal(err) }
    defer echoLn.Close()
    go func() {
        for {
            c, err := echoLn.Accept()
            if err != nil { return }
            go func(c net.Conn) { defer c.Close(); _, _ = io.Copy(c, c) }(c)
        }
    }()

    aptLn, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
        MinVersion:   tls.VersionTLS13,
        Certificates: loadTestCertificate(t, certFile, keyFile),
    })
    if err != nil { t.Fatal(err) }
    srv := New(Config{Addr: aptLn.Addr().String(), Token: "e2e-token"})
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            c, err := aptLn.Accept()
            if err != nil { return }
            go srv.handle(c)
        }
    }()

    sess, err := client.Dial(client.Config{
        ServerAddr:       aptLn.Addr().String(),
        ServerName:       "localhost",
        Token:             "e2e-token",
        InsecureSkipVerify: true,
    })
    if err != nil { t.Fatal(err) }
    defer sess.Close()

    var streams sync.WaitGroup
    for i := 0; i < 8; i++ {
        streams.Add(1)
        go func(i int) {
            defer streams.Done()
            c, err := sess.Open(echoLn.Addr().String())
            if err != nil { t.Error(err); return }
            defer c.Close()
            msg := []byte(fmt.Sprintf("apt-e2e-stream-%d", i))
            if _, err := c.Write(msg); err != nil { t.Error(err); return }
            _ = c.SetReadDeadline(time.Now().Add(3 * time.Second))
            got := make([]byte, len(msg))
            if _, err := io.ReadFull(c, got); err != nil { t.Error(err); return }
            if string(got) != string(msg) { t.Errorf("stream %d got %q want %q", i, got, msg) }
        }(i)
    }
    streams.Wait()

    // Stop the accept loop before waiting for it. The previous test waited
    // first, which could never complete because the listener remained open.
    _ = aptLn.Close()
    wg.Wait()
}

func testCertificate(t *testing.T) (string, string) {
    t.Helper()
    key, err := rsa.GenerateKey(rand.Reader, 2048); if err != nil { t.Fatal(err) }
    tmpl := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "localhost"}, DNSNames: []string{"localhost"}, NotBefore: time.Now().Add(-time.Minute), NotAfter: time.Now().Add(time.Hour), KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}}
    der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key); if err != nil { t.Fatal(err) }
    cf, err := os.CreateTemp("", "apt-e2e-*.crt"); if err != nil { t.Fatal(err) }; defer cf.Close()
    kf, err := os.CreateTemp("", "apt-e2e-*.key"); if err != nil { t.Fatal(err) }; defer kf.Close()
    _ = pem.Encode(cf, &pem.Block{Type: "CERTIFICATE", Bytes: der})
    _ = pem.Encode(kf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
    return cf.Name(), kf.Name()
}

func loadTestCertificate(t *testing.T, certFile, keyFile string) []tls.Certificate {
    t.Helper(); cert, err := tls.LoadX509KeyPair(certFile, keyFile); if err != nil { t.Fatal(err) }; return []tls.Certificate{cert}
}
