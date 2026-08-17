package protocol

import "testing"

func TestReceiveWindowConsumeAndReplenish(t *testing.T) {
    w, err := NewReceiveWindow(100, 200)
    if err != nil { t.Fatal(err) }
    if got := w.Consume(60); got != 60 { t.Fatalf("consume=%d", got) }
    if got := w.Consume(100); got != 40 { t.Fatalf("consume=%d, want 40", got) }
    if w.Available() != 0 { t.Fatalf("available=%d", w.Available()) }
    if err := w.Replenish(100); err != nil { t.Fatal(err) }
    if w.Available() != 100 { t.Fatalf("available=%d", w.Available()) }
}

func TestReceiveWindowRejectsOverflow(t *testing.T) {
    w, err := NewReceiveWindow(100, 100)
    if err != nil { t.Fatal(err) }
    if err := w.Replenish(1); err == nil { t.Fatal("expected overflow") }
}
