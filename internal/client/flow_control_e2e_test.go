package client

import (
    "context"
    "testing"
    "time"
)

// This is the runtime gate for the two-level sender admission path. It uses
// the same BlockingFlowWindow implementation as stream.Write(), but keeps the
// test independent from a real network so it can deterministically prove the
// required BLOCK -> WINDOW_UPDATE -> WAKE transition.
func TestFlowControlE2EBlockUpdateWake(t *testing.T) {
    a, err := newSendAdmission(4, 4)
    if err != nil { t.Fatal(err) }

    if err := a.acquire(context.Background(), 4); err != nil { t.Fatal(err) }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    result := make(chan error, 1)
    go func() { result <- a.acquire(ctx, 1) }()

    select {
    case err := <-result:
        t.Fatalf("sender was not backpressured: %v", err)
    case <-time.After(30 * time.Millisecond):
    }

    if err := a.updateStream(1); err != nil { t.Fatal(err) }
    if err := a.updateSession(1); err != nil { t.Fatal(err) }

    select {
    case err := <-result:
        if err != nil { t.Fatalf("sender did not resume: %v", err) }
    case <-time.After(time.Second):
        t.Fatal("sender remained blocked after both WINDOW_UPDATE events")
    }
}

func TestFlowControlE2ECancelBlockedWrite(t *testing.T) {
    a, err := newSendAdmission(1, 1)
    if err != nil { t.Fatal(err) }
    if err := a.acquire(context.Background(), 1); err != nil { t.Fatal(err) }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
    defer cancel()
    if err := a.acquire(ctx, 1); err == nil {
        t.Fatal("blocked sender unexpectedly succeeded")
    }
}
