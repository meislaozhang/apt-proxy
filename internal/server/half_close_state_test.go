package server

import (
	"testing"
)

func TestHalfCloseStateTransitions(t *testing.T) {
	ss := &serverSession{halfClosed: make(map[uint32]uint8)}
	const id uint32 = 7

	ss.markHalfClosed(id, 1)
	if ss.halfClosed[id] != 1 {
		t.Fatalf("server half-close state = %d, want 1", ss.halfClosed[id])
	}
	if ss.bothHalfClosed(id) {
		t.Fatal("stream must not be fully half-closed after one direction")
	}

	ss.markHalfClosed(id, 2)
	if ss.halfClosed[id] != 3 {
		t.Fatalf("both half-close state = %d, want 3", ss.halfClosed[id])
	}
	if !ss.bothHalfClosed(id) {
		t.Fatal("stream must be fully half-closed after both directions")
	}

	ss.clearHalfClose(id)
	if _, ok := ss.halfClosed[id]; ok {
		t.Fatal("half-close state was not cleared")
	}
}
