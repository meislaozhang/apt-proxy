package protocol

import "testing"

func TestWindowUpdatePayload(t *testing.T) {
    p := WindowUpdatePayload(65536)
    got, err := ParseWindowUpdate(p)
    if err != nil {
        t.Fatal(err)
    }
    if got != 65536 {
        t.Fatalf("got %d, want 65536", got)
    }
}

func TestWindowUpdateRejectsZero(t *testing.T) {
    if _, err := ParseWindowUpdate(WindowUpdatePayload(0)); err == nil {
        t.Fatal("expected zero increment to be rejected")
    }
}

func TestWindowUpdateRejectsBadLength(t *testing.T) {
    if _, err := ParseWindowUpdate([]byte{1, 2, 3}); err == nil {
        t.Fatal("expected malformed payload to be rejected")
    }
}
