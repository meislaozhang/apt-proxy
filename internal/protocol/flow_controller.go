package protocol

import (
    "fmt"
    "sync"
)

// FlowController combines stream and session byte credit.
// A DATA frame may be sent only when both windows have enough credit.
type FlowController struct {
    mu sync.Mutex
    stream *FlowWindow
    session *FlowWindow
}

func NewFlowController(streamInitial, streamMax, sessionInitial, sessionMax uint64) (*FlowController, error) {
    sw, err := NewFlowWindow(streamInitial, streamMax)
    if err != nil { return nil, err }
    cw, err := NewFlowWindow(sessionInitial, sessionMax)
    if err != nil { return nil, err }
    return &FlowController{stream: sw, session: cw}, nil
}

// TryReserve atomically reserves n bytes from both windows, or reserves none.
func (f *FlowController) TryReserve(n uint64) bool {
    if n == 0 { return true }
    f.mu.Lock()
    defer f.mu.Unlock()
    if f.stream.Available() < n || f.session.Available() < n { return false }
    if f.stream.Take(n) != n { return false }
    if f.session.Take(n) != n {
        // This path should be unreachable because both checks and reservation
        // are serialized by f.mu, but keep a defensive rollback.
        _ = f.stream.Add(n)
        return false
    }
    return true
}

func (f *FlowController) UpdateStream(n uint64) error {
    return f.stream.Add(n)
}

func (f *FlowController) UpdateSession(n uint64) error {
    return f.session.Add(n)
}

func (f *FlowController) StreamAvailable() uint64 { return f.stream.Available() }
func (f *FlowController) SessionAvailable() uint64 { return f.session.Available() }

func (f *FlowController) Validate() error {
    if f.stream.Available() > 0 && f.session.Available() > 0 { return nil }
    return fmt.Errorf("flow control blocked")
}
