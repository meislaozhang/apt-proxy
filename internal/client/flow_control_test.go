package client

import (
    "context"
    "testing"
    "time"
)

func TestSendAdmissionBlocksOnStreamCredit(t *testing.T) {
    a, err := newSendAdmission(8, 4)
    if err != nil { t.Fatal(err) }
    if err := a.acquire(context.Background(), 4); err != nil { t.Fatal(err) }

    done := make(chan error, 1)
    go func() { done <- a.acquire(context.Background(), 1) }()

    select {
    case err := <-done:
        t.Fatalf("acquire unexpectedly completed: %v", err)
    case <-time.After(30 * time.Millisecond):
    }

    if err := a.updateStream(1); err != nil { t.Fatal(err) }
    select {
    case err := <-done:
        if err != nil { t.Fatal(err) }
    case <-time.After(time.Second):
        t.Fatal("blocked acquire did not wake after stream WINDOW_UPDATE")
    }
}

func TestSendAdmissionBlocksOnSessionCredit(t *testing.T) {
    a, err := newSendAdmission(4, 8)
    if err != nil { t.Fatal(err) }
    if err := a.acquire(context.Background(), 4); err != nil { t.Fatal(err) }

    done := make(chan error, 1)
    go func() { done <- a.acquire(context.Background(), 1) }()
    select {
    case err := <-done:
        t.Fatalf("acquire unexpectedly completed: %v", err)
    case <-time.After(30 * time.Millisecond):
    }

    if err := a.updateSession(1); err != nil { t.Fatal(err) }
    select {
    case err := <-done:
        if err != nil { t.Fatal(err) }
    case <-time.After(time.Second):
        t.Fatal("blocked acquire did not wake after session WINDOW_UPDATE")
    }
}

func TestSendAdmissionCancellation(t *testing.T) {
    a, err := newSendAdmission(1, 1)
    if err != nil { t.Fatal(err) }
    if err := a.acquire(context.Background(), 1); err != nil { t.Fatal(err) }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
    defer cancel()
    if err := a.acquire(ctx, 1); err == nil { t.Fatal("expected context cancellation") }
}
