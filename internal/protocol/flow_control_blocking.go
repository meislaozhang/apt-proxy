package protocol

import (
    "context"
    "sync"
)

// BlockingFlowWindow provides blocking byte-credit acquisition for DATA
// scheduling. Credit is returned by Add, which wakes blocked senders.
type BlockingFlowWindow struct {
    mu        sync.Mutex
    available uint64
    max       uint64
    notify    chan struct{}
}

func NewBlockingFlowWindow(initial, max uint64) (*BlockingFlowWindow, error) {
    w, err := NewFlowWindow(initial, max)
    if err != nil { return nil, err }
    return &BlockingFlowWindow{available: w.available, max: w.max, notify: make(chan struct{})}, nil
}

func (w *BlockingFlowWindow) Available() uint64 {
    w.mu.Lock(); defer w.mu.Unlock()
    return w.available
}

func (w *BlockingFlowWindow) Acquire(ctx context.Context, n uint64) error {
    if n == 0 { return nil }
    for {
        w.mu.Lock()
        if w.available >= n {
            w.available -= n
            w.mu.Unlock()
            return nil
        }
        ch := w.notify
        w.mu.Unlock()
        select {
        case <-ctx.Done(): return ctx.Err()
        case <-ch:
        }
    }
}

func (w *BlockingFlowWindow) Add(n uint64) error {
    if n == 0 { return nil }
    w.mu.Lock()
    if n > w.max-w.available { w.mu.Unlock(); return ErrFlowWindowOverflow }
    w.available += n
    old := w.notify
    w.notify = make(chan struct{})
    close(old)
    w.mu.Unlock()
    return nil
}
