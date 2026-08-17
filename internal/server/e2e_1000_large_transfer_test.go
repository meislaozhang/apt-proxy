package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

// TestE2E1000StreamLargeTransferIntegrity is the final multiplexing/data-integrity
// stress gate. It is opt-in because it creates 1000 concurrent streams and moves
// 64 MiB through one multiplexed session.
func TestE2E1000StreamLargeTransferIntegrity(t *testing.T) {
	if os.Getenv("APT_PROXY_1000_LARGE_E2E") != "1" {
		t.Skip("set APT_PROXY_1000_LARGE_E2E=1 for the 1000-stream large-transfer gate")
	}

	const streams = 1000
	const size = 64 << 10 // 64 KiB per stream; 64 MiB total.

	echoLn, sess, cleanup := testE2ESession(t)
	defer cleanup()

	var wg sync.WaitGroup
	errCh := make(chan error, streams)
	start := time.Now()

	for streamID := 0; streamID < streams; streamID++ {
		streamID := streamID
		wg.Add(1)
		go func() {
			defer wg.Done()
			c, err := sess.Open(echoLn.Addr().String())
			if err != nil {
				errCh <- fmt.Errorf("stream %d open: %w", streamID, err)
				return
			}
			defer c.Close()

			payload := make([]byte, size)
			for i := range payload {
				payload[i] = byte((i*31 + streamID*17) % 251)
			}
			want := sha256.Sum256(payload)

			if _, err := io.Copy(c, bytes.NewReader(payload)); err != nil {
				errCh <- fmt.Errorf("stream %d write: %w", streamID, err)
				return
			}

			_ = c.SetReadDeadline(time.Now().Add(60 * time.Second))
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
		}()
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Error(err)
	}
	t.Logf("1000-stream integrity gate completed in %s", time.Since(start).Round(time.Millisecond))
}
