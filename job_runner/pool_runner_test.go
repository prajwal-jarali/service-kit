package job_runner

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPool_BasicExecution(t *testing.T) {
	p := NewPool(3, 10)
	defer p.Stop()

	var count atomic.Int32

	for range 5 {
		p.Submit(func(ctx context.Context) {
			count.Add(1)
		})
	}

	time.Sleep(100 * time.Millisecond)

	if count.Load() != 5 {
		t.Fatalf("expected 5 tasks executed")
	}
}

func TestPool_ConcurrencyLimit(t *testing.T) {
	p := NewPool(2, 10)
	defer p.Stop()

	var running atomic.Int32
	var maxRunning int32
	var mu sync.Mutex

	for range 5 {
		p.Submit(func(ctx context.Context) {
			r := running.Add(1)

			mu.Lock()
			if r > maxRunning {
				maxRunning = r
			}
			mu.Unlock()

			time.Sleep(50 * time.Millisecond)

			running.Add(-1)
		})
	}

	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	finalMax := maxRunning
	mu.Unlock()

	if finalMax > 2 {
		t.Fatalf("expected max concurrency 2, got %d", finalMax)
	}
}

func TestPool_Stop(t *testing.T) {
	p := NewPool(2, 10)

	var count atomic.Int32

	for range 5 {
		p.Submit(func(ctx context.Context) {
			time.Sleep(50 * time.Millisecond)
			count.Add(1)
		})
	}

	p.Stop()

	if count.Load() == 0 {
		t.Fatalf("tasks should complete before stop")
	}
}
