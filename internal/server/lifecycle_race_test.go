package server

import (
	"net"
	"sync"
	"testing"

	"github.com/meislaozhang/apt-proxy/internal/protocol"
)

func TestLifecycleResetCloseHalfCloseConcurrent(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	ss := &serverSession{
		c: c1,
		streams: map[uint32]net.Conn{1: c2},
		streamSend: make(map[uint32]*protocol.BlockingFlowWindow),
		halfClosed: map[uint32]uint8{1: 1},
	}

	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(3)
		go func() { defer wg.Done(); ss.reset(1) }()
		go func() { defer wg.Done(); ss.closeStream(1) }()
		go func() { defer wg.Done(); ss.markHalfClosed(1, 2) }()
	}
	wg.Wait()

	ss.mu.Lock()
	defer ss.mu.Unlock()
	if _, ok := ss.streams[1]; ok {
		t.Fatal("stream remains after lifecycle race")
	}
	if _, ok := ss.halfClosed[1]; ok {
		t.Fatal("half-close state remains after lifecycle race")
	}
}
