package protocol

import (
    "fmt"
    "sync"
)

// ReceiveWindow tracks how much DATA the peer may still send before credit
// must be refreshed. Consumed bytes are returned as WINDOW_UPDATE credit by
// the caller after the application has accepted them.
type ReceiveWindow struct {
    mu        sync.Mutex
    available uint64
    max       uint64
}

func NewReceiveWindow(initial, max uint64) (*ReceiveWindow, error) {
    if initial == 0 || max == 0 || initial > max {
        return nil, fmt.Errorf("invalid receive window initial=%d max=%d", initial, max)
    }
    return &ReceiveWindow{available: initial, max: max}, nil
}

func (w *ReceiveWindow) Available() uint64 {
    w.mu.Lock()
    defer w.mu.Unlock()
    return w.available
}

// Consume admits up to n incoming bytes and returns the amount accepted.
func (w *ReceiveWindow) Consume(n uint64) uint64 {
    w.mu.Lock()
    defer w.mu.Unlock()
    if n > w.available {
        n = w.available
    }
    w.available -= n
    return n
}

func (w *ReceiveWindow) Replenish(n uint64) error {
    if n == 0 {
        return fmt.Errorf("receive-window increment must be non-zero")
    }
    w.mu.Lock()
    defer w.mu.Unlock()
    if n > w.max-w.available {
        return fmt.Errorf("receive-window update exceeds maximum")
    }
    w.available += n
    return nil
}
