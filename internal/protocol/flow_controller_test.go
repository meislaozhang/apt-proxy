package protocol

import "testing"

func TestFlowControllerRequiresBothWindows(t *testing.T) {
    f, err := NewFlowController(10, 100, 20, 100)
    if err != nil { t.Fatal(err) }
    if !f.TryReserve(10) { t.Fatal("expected reservation") }
    if f.StreamAvailable() != 0 || f.SessionAvailable() != 10 { t.Fatalf("unexpected windows: stream=%d session=%d", f.StreamAvailable(), f.SessionAvailable()) }
    if f.TryReserve(1) { t.Fatal("expected stream window to block") }
    if err := f.UpdateStream(10); err != nil { t.Fatal(err) }
    if !f.TryReserve(11) { t.Fatal("expected session window to block only when request exceeds remaining credit") }
}

func TestFlowControllerSessionBlocks(t *testing.T) {
    f, err := NewFlowController(100, 100, 5, 100)
    if err != nil { t.Fatal(err) }
    if f.TryReserve(6) { t.Fatal("expected session window to block") }
    if !f.TryReserve(5) { t.Fatal("expected five bytes to fit") }
}
