package protocol

import (
    "fmt"
    "sync"
)

// SendWindow limits DATA bytes that a stream may have in flight.
// It is deliberately independent from transport framing so it can later be
// reused by the TCP and QUIC transports.
type SendWindow struct {
    mu       sync.Mutex
    available uint64
    max       uint64
}

func NewSendWindow(initial, max uint64) (*SendWindow, error) {
    if initial == 0 || max == 0 || initial > max {
        return nil, fmt.Errorf("invalid send window initial=%d max=%d", initial, max)
    }
    return &SendWindow{available: initial, max: max}, nil
}

func (w *SendWindow) Available() uint64 {
    w.mu.Lock()
    defer w.mu.Unlock()
    return w.available
}

// Acquire reserves up to n bytes and returns the amount granted.
func (w *SendWindow) Acquire(n uint64) uint64 {
    w.mu.Lock()
    defer w.mu.Unlock()
    if n > w.available {
        n = w.available
    }
    w.available -= n
    return n
}

// Update adds credit while preventing uint64 overflow and max-window growth.
func (w *SendWindow) Update(n uint64) error {
    if n == 0 {
        return fmt.Errorf("window increment must be non-zero")
    }
    w.mu.Lock()
    defer w.mu.Unlock()
    if n > w.max-w.available {
        return fmt.Errorf("window update exceeds maximum")
    }
    w.available += n
    return nil
}
