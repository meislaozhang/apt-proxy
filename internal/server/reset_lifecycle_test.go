package server

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/meislaozhang/apt-proxy/internal/protocol"
)

func TestResetClearsStreamLifecycleState(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	originConn, originPeer := net.Pipe()
	defer originPeer.Close()

	ss := &serverSession{
		c:          serverConn,
		r:          bufio.NewReader(serverConn),
		streams:    map[uint32]net.Conn{7: originConn},
		streamSend: make(map[uint32]*protocol.BlockingFlowWindow),
		halfClosed: map[uint32]uint8{7: 3},
	}

	// Arm the peer reader before reset starts. net.Pipe writes are synchronous;
	// starting the reader first removes a scheduling race in this lifecycle test.
	frameCh := make(chan struct {
		frame protocol.Frame
		err   error
	}, 1)
	go func() {
		f, err := protocol.ReadFrame(bufio.NewReader(clientConn))
		frameCh <- struct {
			frame protocol.Frame
			err   error
		}{f, err}
	}()

	resetDone := make(chan struct{})
	go func() {
		ss.reset(7)
		close(resetDone)
	}()

	// Do not race resetDone against frameCh: a successful synchronous Write may
	// unblock and let reset return before the reader goroutine gets scheduled to
	// publish its already-received frame. The protocol event is the assertion;
	// wait for it with a bounded timeout, then wait for reset to finish.
	select {
	case result := <-frameCh:
		if result.err != nil {
			t.Fatalf("read reset frame: %v", result.err)
		}
		if result.frame.Type != protocol.TypeReset || result.frame.StreamID != 7 {
			t.Fatalf("got frame type=%v stream=%d, want RESET stream=7", result.frame.Type, result.frame.StreamID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for RESET frame")
	}

	select {
	case <-resetDone:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for reset to finish")
	}

	ss.mu.Lock()
	_, streamExists := ss.streams[7]
	_, windowExists := ss.streamSend[7]
	_, halfCloseExists := ss.halfClosed[7]
	ss.mu.Unlock()
	if streamExists || windowExists || halfCloseExists {
		t.Fatal("reset did not clear stream lifecycle state")
	}

	if _, err := originPeer.Write([]byte("x")); err == nil {
		t.Fatal("origin peer write unexpectedly succeeded after reset")
	}

	_ = serverConn.Close()
}
