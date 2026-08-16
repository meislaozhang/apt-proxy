package protocol

import (
    "bufio"
    "bytes"
    "testing"
)

func TestFrameRoundTrip(t *testing.T) {
    want := Frame{Type: TypeData, Flags: 3, StreamID: 42, Payload: []byte("hello apt")}
    var b bytes.Buffer
    if err := want.Write(&b); err != nil { t.Fatal(err) }
    got, err := ReadFrame(bufio.NewReader(&b))
    if err != nil { t.Fatal(err) }
    if got.Type != want.Type || got.Flags != want.Flags || got.StreamID != want.StreamID || !bytes.Equal(got.Payload, want.Payload) { t.Fatalf("got %#v want %#v", got, want) }
}

func TestFrameRoundTripMultipleStreams(t *testing.T) {
    var b bytes.Buffer
    frames := []Frame{
        {Type: TypeOpen, StreamID: 1, Payload: []byte("example.com:443")},
        {Type: TypeData, StreamID: 3, Payload: []byte("stream-three")},
        {Type: TypeData, StreamID: 1, Payload: []byte("stream-one")},
        {Type: TypeClose, StreamID: 3},
    }
    for _, f := range frames { if err := f.Write(&b); err != nil { t.Fatal(err) } }
    r := bufio.NewReader(&b)
    for i, want := range frames {
        got, err := ReadFrame(r)
        if err != nil { t.Fatalf("frame %d: %v", i, err) }
        if got.Type != want.Type || got.StreamID != want.StreamID || !bytes.Equal(got.Payload, want.Payload) { t.Fatalf("frame %d: got %#v want %#v", i, got, want) }
    }
}

func TestFrameRejectsOversizedPayload(t *testing.T) {
    f := Frame{Type: TypeData, StreamID: 1, Payload: make([]byte, MaxPayload+1)}
    var b bytes.Buffer
    if err := f.Write(&b); err == nil { t.Fatal("expected oversized payload to be rejected") }
}
