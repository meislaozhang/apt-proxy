package server

import (
	"net"
	"sync"
	"testing"

	"github.com/meislaozhang/apt-proxy/internal/protocol"
)

func TestMultiStreamLifecycleCleanup(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	const streams = 100
	ss := &serverSession{
		c: c1,
		streams: make(map[uint32]net.Conn, streams),
		streamSend: make(map[uint32]*protocol.BlockingFlowWindow, streams),
		halfClosed: make(map[uint32]uint8, streams),
	}
	for i := 1; i <= streams; i++ {
		ss.streams[uint32(i)] = c2
		ss.streamSend[uint32(i)] = protocol.NewBlockingFlowWindow(1)
		ss.halfClosed[uint32(i)] = 1
	}

	var wg sync.WaitGroup
	for i := 1; i <= streams; i++ {
		id := uint32(i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			ss.reset(id)
		}()
	}
	wg.Wait()

	ss.mu.Lock()
	defer ss.mu.Unlock()
	if len(ss.streams) != 0 || len(ss.streamSend) != 0 || len(ss.halfClosed) != 0 {
		t.Fatalf("lifecycle state leaked: streams=%d sendWindows=%d halfClosed=%d", len(ss.streams), len(ss.streamSend), len(ss.halfClosed))
	}
}
