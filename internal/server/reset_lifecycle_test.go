package server

import (
	"bufio"
	"net"
	"testing"

	"github.com/meislaozhang/apt-proxy/internal/protocol"
)

func TestResetClearsStreamLifecycleState(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	originConn, originPeer := net.Pipe()
	defer originPeer.Close()

	ss := &serverSession{
		c:           serverConn,
		r:           bufio.NewReader(serverConn),
		streams:     map[uint32]net.Conn{7: originConn},
		streamSend:  make(map[uint32]*protocol.BlockingFlowWindow),
		halfClosed:  map[uint32]uint8{7: 3},
	}

	ss.reset(7)

	ss.mu.Lock()
	_, streamExists := ss.streams[7]
	_, windowExists := ss.streamSend[7]
	_, halfCloseExists := ss.halfClosed[7]
	ss.mu.Unlock()
	if streamExists || windowExists || halfCloseExists {
		t.Fatal("reset did not clear stream lifecycle state")
	}

	f, err := protocol.ReadFrame(bufio.NewReader(clientConn))
	if err != nil {
		t.Fatalf("read reset frame: %v", err)
	}
	if f.Type != protocol.TypeReset || f.StreamID != 7 {
		t.Fatalf("got frame type=%v stream=%d, want RESET stream=7", f.Type, f.StreamID)
	}

	if _, err := originPeer.Write([]byte("x")); err == nil {
		t.Fatal("origin peer write unexpectedly succeeded after reset")
	}

	_ = serverConn.Close()
}
