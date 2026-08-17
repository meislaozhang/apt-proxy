package server

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/meislaozhang/apt-proxy/internal/client"
)

// TestE2ELargeTransferIntegrity verifies byte-for-byte integrity through the
// multiplexed TCP/TLS path. The default sizes stay CI-friendly; setting
// APT_PROXY_LARGE_E2E=1 adds a 100 MiB case, while
// APT_PROXY_LARGE_E2E_GIB=1 adds a 1 GiB case for dedicated runners.
func TestE2ELargeTransferIntegrity(t *testing.T) {
	sizes := []int{1 << 20, 16 << 20, 64 << 20}
	if os.Getenv("APT_PROXY_LARGE_E2E") == "1" {
		sizes = append(sizes, 100<<20)
	}
	if os.Getenv("APT_PROXY_LARGE_E2E_GIB") == "1" {
		sizes = append(sizes, 1<<30)
	}

	echoLn, sess, cleanup := testE2ESession(t)
	defer cleanup()

	for _, size := range sizes {
		size := size
		t.Run(fmt.Sprintf("bytes-%d", size), func(t *testing.T) {
			c, err := sess.Open(echoLn.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()
			_ = c.SetReadDeadline(time.Now().Add(30 * time.Second))

			payload := make([]byte, size)
			for i := range payload {
				payload[i] = byte((i*31 + 17) % 251)
			}
			want := sha256.Sum256(payload)

			if _, err := io.Copy(c, bytes.NewReader(payload)); err != nil {
				t.Fatalf("write %d bytes: %v", size, err)
			}

			gotHash := sha256.New()
			if _, err := io.Copy(gotHash, io.LimitReader(c, int64(size))); err != nil {
				t.Fatalf("read %d bytes: %v", size, err)
			}
			got := gotHash.Sum(nil)
			if !bytes.Equal(got, want[:]) {
				t.Fatalf("sha256 mismatch: got %s want %s", hex.EncodeToString(got), hex.EncodeToString(want[:]))
			}
		})
	}
}

func TestE2EConcurrentLargeTransferIntegrity(t *testing.T) {
	if os.Getenv("APT_PROXY_LARGE_E2E") != "1" {
		t.Skip("set APT_PROXY_LARGE_E2E=1 for the concurrent large-transfer gate")
	}

	const streams = 32
	const size = 1 << 20
	echoLn, sess, cleanup := testE2ESession(t)
	defer cleanup()

	errCh := make(chan error, streams)
	for streamID := 0; streamID < streams; streamID++ {
		streamID := streamID
		go func() {
			c, err := sess.Open(echoLn.Addr().String())
			if err != nil {
				errCh <- fmt.Errorf("stream %d open: %w", streamID, err)
				return
			}
			defer c.Close()
			_ = c.SetReadDeadline(time.Now().Add(30 * time.Second))

			payload := make([]byte, size)
			for i := range payload {
				payload[i] = byte((i*17 + streamID*13) % 251)
			}
			want := sha256.Sum256(payload)
			if _, err := io.Copy(c, bytes.NewReader(payload)); err != nil {
				errCh <- fmt.Errorf("stream %d write: %w", streamID, err)
				return
			}
			gotHash := sha256.New()
			if _, err := io.Copy(gotHash, io.LimitReader(c, int64(size))); err != nil {
				errCh <- fmt.Errorf("stream %d read: %w", streamID, err)
				return
			}
			got := gotHash.Sum(nil)
			if !bytes.Equal(got, want[:]) {
				errCh <- fmt.Errorf("stream %d sha256 mismatch: got %s want %s", streamID, hex.EncodeToString(got), hex.EncodeToString(want[:]))
				return
			}
			errCh <- nil
		}()
	}
	for i := 0; i < streams; i++ {
		if err := <-errCh; err != nil {
			t.Error(err)
		}
	}
}

func testE2ESession(t *testing.T) (net.Listener, *client.Session, func()) {
	t.Helper()
	certFile, keyFile := testCertificate(t)

	echoLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			c, err := echoLn.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				_, _ = io.Copy(c, c)
			}(c)
		}
	}()

	aptLn, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: loadTestCertificate(t, certFile, keyFile),
	})
	if err != nil {
		echoLn.Close()
		t.Fatal(err)
	}
	srv := New(Config{Addr: aptLn.Addr().String(), Token: "e2e-token"})
	go func() {
		for {
			c, err := aptLn.Accept()
			if err != nil {
				return
			}
			go srv.handle(c)
		}
	}()

	sess, err := client.Dial(client.Config{
		ServerAddr:         aptLn.Addr().String(),
		ServerName:         "localhost",
		Token:              "e2e-token",
		InsecureSkipVerify: true,
	})
	if err != nil {
		aptLn.Close()
		echoLn.Close()
		os.Remove(certFile)
		os.Remove(keyFile)
		t.Fatal(err)
	}

	return echoLn, sess, func() {
		sess.Close()
		aptLn.Close()
		echoLn.Close()
		os.Remove(certFile)
		os.Remove(keyFile)
	}
}
