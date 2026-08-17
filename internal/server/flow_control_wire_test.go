package server

import (
    "bufio"
    "net"
    "testing"

    "github.com/meislaozhang/apt-proxy/internal/protocol"
)

// TestWriteDataWireFlowControl verifies the server's real frame writer emits
// DATA followed by stream- and session-level WINDOW_UPDATE frames. net.Pipe
// gives us an actual byte stream, so this exercises Frame.Write/ReadFrame and
// the same write path used by a live TLS session without requiring a public VPS.
func TestWriteDataWireFlowControl(t *testing.T) {
    wire, peer := net.Pipe()
    defer wire.Close()
    defer peer.Close()

    ss := &serverSession{
        c: wire,
        streams: map[uint32]net.Conn{},
    }
    origin, originPeer := net.Pipe()
    defer origin.Close()
    defer originPeer.Close()
    const streamID uint32 = 7
    ss.streams[streamID] = origin

    payload := []byte("wire-flow-control")
    done := make(chan struct{})
    go func() {
        ss.writeData(streamID, payload)
        close(done)
    }()

    r := bufio.NewReader(peer)
    f1, err := protocol.ReadFrame(r)
    if err != nil { t.Fatal(err) }
    if f1.Type != protocol.TypeData || f1.StreamID != streamID || string(f1.Payload) != string(payload) {
        t.Fatalf("unexpected DATA frame: %+v", f1)
    }
    f2, err := protocol.ReadFrame(r)
    if err != nil { t.Fatal(err) }
    if f2.Type != protocol.TypeWindowUpdate || f2.StreamID != streamID {
        t.Fatalf("unexpected stream WINDOW_UPDATE: %+v", f2)
    }
    f3, err := protocol.ReadFrame(r)
    if err != nil { t.Fatal(err) }
    if f3.Type != protocol.TypeWindowUpdate || f3.StreamID != 0 {
        t.Fatalf("unexpected session WINDOW_UPDATE: %+v", f3)
    }
    if got := string(readExact(t, originPeer, len(payload))); got != string(payload) {
        t.Fatalf("origin received %q, want %q", got, payload)
    }
    <-done
}

func readExact(t *testing.T, c net.Conn, n int) []byte {
    t.Helper()
    b := make([]byte, n)
    if _, err := c.Read(b); err != nil { t.Fatal(err) }
    return b
}
