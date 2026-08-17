package client

import (
    "context"
    "fmt"

    "github.com/meislaozhang/apt-proxy/internal/protocol"
)

// sendAdmission enforces both per-stream and per-session byte credits before
// a DATA frame is put on the wire. Credits are acquired before the session
// write lock so a blocked stream cannot serialize the whole session.
type sendAdmission struct {
    session *protocol.BlockingFlowWindow
    stream  *protocol.BlockingFlowWindow
}

func newSendAdmission(sessionWindow, streamWindow uint64) (*sendAdmission, error) {
    sw, err := protocol.NewBlockingFlowWindow(sessionWindow, sessionWindow)
    if err != nil { return nil, err }
    tw, err := protocol.NewBlockingFlowWindow(streamWindow, streamWindow)
    if err != nil { return nil, err }
    return &sendAdmission{session: sw, stream: tw}, nil
}

func (a *sendAdmission) acquire(ctx context.Context, n uint64) error {
    if n == 0 { return nil }
    if err := a.stream.Acquire(ctx, n); err != nil { return err }
    if err := a.session.Acquire(ctx, n); err != nil {
        if rerr := a.stream.Add(n); rerr != nil {
            return fmt.Errorf("restore stream credit: %w (original: %v)", rerr, err)
        }
        return err
    }
    return nil
}

func (a *sendAdmission) updateSession(n uint64) error { return a.session.Add(n) }
func (a *sendAdmission) updateStream(n uint64) error { return a.stream.Add(n) }
