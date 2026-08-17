package server

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/meislaozhang/apt-proxy/internal/protocol"
)

// TestPumpWaitsForSessionWindow proves that an exhausted session window
// backpressures Server -> Client DATA instead of closing the stream. The
// blocked pump must resume after a session WINDOW_UPDATE replenishes credit.
func TestPumpWaitsForSessionWindow(t *testing.T) {
	wire, peer := net.Pipe()
	defer wire.Close()
	defer peer.Close()

	origin, originPeer := net.Pipe()
	defer origin.Close()
	defer originPeer.Close()

	// The window implementation rejects an initial credit of zero. Seed one
	// byte, consume it, and thereby put the pump into the intended zero-credit
	// state before the payload arrives.
	sessionWindow, err := protocol.NewBlockingFlowWindow(1, 64*1024)
	if err != nil {
		t.Fatal(err)
	}
	if err := sessionWindow.Acquire(nilContext{}, 1); err != nil {
		t.Fatal(err)
	}
	streamWindow, err := protocol.NewBlockingFlowWindow(64*1024, 64*1024)
	if err != nil {
		t.Fatal(err)
	}

	const streamID uint32 = 9
	ss := &serverSession{
		c: wire,
		streams: map[uint32]net.Conn{streamID: origin},
		send: sessionWindow,
		streamSend: map[uint32]*protocol.BlockingFlowWindow{streamID: streamWindow},
	}

	pumpDone := make(chan struct{})
	go func() {
		ss.pump(streamID, origin, streamWindow)
		close(pumpDone)
	}()

	payload := []byte("blocked-until-window-update")
	writeDone := make(chan error, 1)
	go func() {
		_, err := originPeer.Write(payload)
		writeDone <- err
	}()

	// The pump can read the origin payload, but must remain blocked because
	// the session window has zero credit. No DATA frame may reach the wire.
	time.Sleep(50 * time.Millisecond)
	r := bufio.NewReader(peer)
	_ = peer.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
	_, err = r.ReadByte()
	if err == nil {
		t.Fatal("DATA reached wire before session WINDOW_UPDATE")
	}
	_ = peer.SetReadDeadline(time.Time{})

	if err := sessionWindow.Add(uint64(len(payload))); err != nil {
		t.Fatal(err)
	}

	f, err := protocol.ReadFrame(r)
	if err != nil {
		t.Fatal(err)
	}
	if f.Type != protocol.TypeData || f.StreamID != streamID || string(f.Payload) != string(payload) {
		t.Fatalf("unexpected resumed DATA frame: %+v", f)
	}

	select {
	case err := <-writeDone:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("origin write did not complete")
	}

	_ = originPeer.Close()
	select {
	case <-pumpDone:
	case <-time.After(time.Second):
		t.Fatal("pump did not exit after origin close")
	}
}

type nilContext struct{}
func (nilContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (nilContext) Done() <-chan struct{} { return nil }
func (nilContext) Err() error { return nil }
func (nilContext) Value(any) any { return nil }
