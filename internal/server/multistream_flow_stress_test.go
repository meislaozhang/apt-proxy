package server

import (
	"sync"
	"testing"

	"github.com/meislaozhang/apt-proxy/internal/protocol"
)

func Test100StreamsConcurrentFlowWindows(t *testing.T) {
	const streams = 100
	const rounds = 32

	windows := make([]*protocol.BlockingFlowWindow, streams)
	for i := range windows {
		w, err := protocol.NewBlockingFlowWindow(64*1024, 64*1024)
		if err != nil {
			t.Fatalf("stream %d window: %v", i+1, err)
		}
		windows[i] = w
	}

	var wg sync.WaitGroup
	for i := 0; i < streams; i++ {
		w := windows[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := 0; r < rounds; r++ {
				const n = uint64(1024)
				if err := w.Acquire(nilContext{}, n); err != nil {
					t.Errorf("Acquire: %v", err)
					return
				}
				if err := w.Add(n); err != nil {
					t.Errorf("Add: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()
}

// nilContext is intentionally tiny: BlockingFlowWindow only requires a
// context-like Done/Err contract for this stress gate.
type nilContext struct{}
func (nilContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (nilContext) Done() <-chan struct{} { return nil }
func (nilContext) Err() error { return nil }
func (nilContext) Value(any) any { return nil }
