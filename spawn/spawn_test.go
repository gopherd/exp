package spawn_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gopherd/exp/spawn"
)

func TestRun(t *testing.T) {
	ctx := context.Background()
	var called int32

	handle := spawn.Run(ctx, func(ctx context.Context) {
		atomic.AddInt32(&called, 1)
	})

	// Wait for the task to complete
	handle.Join(ctx)

	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("Expected function to be called once, got %d", called)
	}
}

func TestRun_Cancel(t *testing.T) {
	ctx := context.Background()
	var called int32

	handle := spawn.Run(ctx, func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				atomic.AddInt32(&called, 1)
				time.Sleep(10 * time.Millisecond)
			}
		}
	})

	// Let the task run for a short while
	time.Sleep(50 * time.Millisecond)

	// Cancel the task
	handle.Cancel()

	// Wait for the task to finish
	handle.Join(ctx)

	if atomic.LoadInt32(&called) == 0 {
		t.Errorf("Expected function to be called at least once")
	}
}

func TestTick(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var called int32

	handle := spawn.Tick(ctx, func(ctx context.Context) {
		atomic.AddInt32(&called, 1)
	}, 50*time.Millisecond)

	// Wait for the task to complete
	handle.Join(ctx)

	count := atomic.LoadInt32(&called)
	if count < 3 || count > 5 {
		t.Errorf("Expected function to be called between 3 and 5 times, got %d", count)
	}
}

func TestTick_Cancel(t *testing.T) {
	ctx := context.Background()
	var called int32

	handle := spawn.Tick(ctx, func(ctx context.Context) {
		atomic.AddInt32(&called, 1)
	}, 50*time.Millisecond)

	// Let the task run for a short while
	time.Sleep(120 * time.Millisecond)

	// Cancel the task
	handle.Cancel()

	// Wait for the task to finish
	handle.Join(ctx)

	count := atomic.LoadInt32(&called)
	if count < 2 {
		t.Errorf("Expected function to be called at least 2 times, got %d", count)
	}
}
