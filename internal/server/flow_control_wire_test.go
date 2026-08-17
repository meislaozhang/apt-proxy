package server

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/meislaozhang/apt-proxy/internal/protocol"
)

// TestWriteDataWireFlowControl verifies the server's real frame writer emits
// DATA-window accounting updates at both stream and session scope. net.Pipe
// gives us an actual byte stream, so this exercises Frame.Write/ReadFrame and
// the same write path used by a live TLS session without requiring a public VPS.
func TestWriteDataWireFlowControl(t *testing.T) {
	wire, peer := net.Pipe()
	defer wire.Close()
	defer peer.Close()

	ss := &serverSession{
		c:       wire,
		streams: map[uint32]net.Conn{},
	}
	origin, originPeer := net.Pipe()
	defer origin.Close()
	defer originPeer.Close()
	const streamID uint32 = 7
	ss.streams[streamID] = origin

	payload := []byte("wire-flow-control")
	originRead := make(chan error, 1)
	go func() {
		buf := make([]byte, len(payload))
		_, err := originPeer.Read(buf)
		if err == nil && string(buf) != string(payload) {
			originRead <- &payloadMismatch{got: string(buf), want: string(payload)}
			return
		}
		originRead <- err
	}()

	done := make(chan struct{})
	go func() {
		ss.writeData(streamID, payload)
		close(done)
	}()

	r := bufio.NewReader(peer)
	f1, err := protocol.ReadFrame(r)
	if err != nil {
		t.Fatal(err)
	}
	if f1.Type != protocol.TypeWindowUpdate || f1.StreamID != streamID {
		t.Fatalf("unexpected first frame: %+v", f1)
	}
	f2, err := protocol.ReadFrame(r)
	if err != nil {
		t.Fatal(err)
	}
	if f2.Type != protocol.TypeWindowUpdate || f2.StreamID != 0 {
		t.Fatalf("unexpected second frame: %+v", f2)
	}

	select {
	case err := <-originRead:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for origin write")
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for writeData to finish")
	}
}

type payloadMismatch struct{ got, want string }

func (e *payloadMismatch) Error() string {
	return "origin payload mismatch: got " + e.got + " want " + e.want
}
