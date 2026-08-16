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
    if got.Type != want.Type || got.Flags != want.Flags || got.StreamID != want.StreamID || !bytes.Equal(got.Payload, want.Payload) {
        t.Fatalf("got %#v want %#v", got, want)
    }
}
