package protocol

import "testing"

func TestWindowUpdateRoundTrip(t *testing.T) {
    for _, want := range []uint64{1, 4096, 1 << 20, ^uint64(0)} {
        got, ok := DecodeWindowUpdate(EncodeWindowUpdate(want))
        if !ok || got != want { t.Fatalf("got=%d ok=%v want=%d", got, ok, want) }
    }
}

func TestWindowUpdateRejectsInvalidPayload(t *testing.T) {
    if _, ok := DecodeWindowUpdate([]byte{1}); ok { t.Fatal("accepted short window update") }
    if _, ok := DecodeWindowUpdate(make([]byte, 8)); ok { t.Fatal("accepted zero increment") }
}
