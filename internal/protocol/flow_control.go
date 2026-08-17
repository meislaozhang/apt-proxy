package protocol

import (
    "fmt"
    "sync"
)

// FlowWindow is a bounded byte-credit counter used by a stream or session.
// A zero credit means the sender must apply backpressure until credit is added.
type FlowWindow struct {
    mu        sync.Mutex
    available uint64
    max       uint64
}

func NewFlowWindow(initial, max uint64) (*FlowWindow, error) {
    if initial == 0 || max == 0 || initial > max {
        return nil, fmt.Errorf("invalid flow window initial=%d max=%d", initial, max)
    }
    return &FlowWindow{available: initial, max: max}, nil
}

func (w *FlowWindow) Available() uint64 {
    w.mu.Lock()
    defer w.mu.Unlock()
    return w.available
}

func (w *FlowWindow) Take(n uint64) uint64 {
    w.mu.Lock()
    defer w.mu.Unlock()
    if n > w.available {
        n = w.available
    }
    w.available -= n
    return n
}

func (w *FlowWindow) Add(n uint64) error {
    if n == 0 {
        return fmt.Errorf("flow-window increment must be non-zero")
    }
    w.mu.Lock()
    defer w.mu.Unlock()
    if n > w.max-w.available {
        return fmt.Errorf("flow-window update exceeds maximum")
    }
    w.available += n
    return nil
}

// TryTake is a non-blocking admission check for DATA scheduling.
func (w *FlowWindow) TryTake(n uint64) (uint64, bool) {
    got := w.Take(n)
    return got, got == n
}
