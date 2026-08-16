package protocol

import (
    "testing"
)

func TestSendWindowAcquireAndUpdate(t *testing.T) {
    w, err := NewSendWindow(100, 200)
    if err != nil { t.Fatal(err) }
    if got := w.Acquire(60); got != 60 { t.Fatalf("acquire=%d", got) }
    if got := w.Acquire(100); got != 40 { t.Fatalf("acquire=%d, want 40", got) }
    if w.Available() != 0 { t.Fatalf("available=%d", w.Available()) }
    if err := w.Update(100); err != nil { t.Fatal(err) }
    if w.Available() != 100 { t.Fatalf("available=%d", w.Available()) }
}

func TestSendWindowRejectsOverflow(t *testing.T) {
    w, err := NewSendWindow(100, 100)
    if err != nil { t.Fatal(err) }
    if err := w.Update(1); err == nil { t.Fatal("expected overflow") }
}

func TestSendWindowRejectsInvalidConstruction(t *testing.T) {
    if _, err := NewSendWindow(0, 1); err == nil { t.Fatal("expected invalid initial") }
    if _, err := NewSendWindow(2, 1); err == nil { t.Fatal("expected initial > max") }
}
