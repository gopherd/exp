// Package spawn provides functionality to run and manage concurrent tasks.
package spawn

import (
	"context"
	"time"
)

// Handle defines methods to control concurrent tasks.
type Handle interface {
	// Join waits for the task to complete or the context to be canceled.
	Join(context.Context)
	// Cancel stops the execution of the task.
	Cancel()
}

// taskHandle implements the Handle interface and contains control information for a task.
type taskHandle struct {
	done   chan struct{}
	cancel context.CancelFunc
}

// Join blocks until the task completes or the context is canceled.
func (h *taskHandle) Join(ctx context.Context) {
	select {
	case <-ctx.Done():
		// Context canceled or timed out
	case <-h.done:
		// Task completed
	}
}

// Cancel stops the execution of the task.
func (h *taskHandle) Cancel() {
	if h.cancel != nil {
		h.cancel()
	}
}

// Run starts a new concurrent task with the given context and function.
//
// Parameters:
//   - ctx: The context used to control the lifecycle of the task.
//   - fn:  The task function that accepts a context.
//
// Returns:
//   - Handle: A handle that can be used to control the task.
func Run(ctx context.Context, f func(context.Context)) Handle {
	ctx, cancel := context.WithCancel(ctx)
	h := &taskHandle{
		done:   make(chan struct{}),
		cancel: cancel,
	}
	go func() {
		defer close(h.done)
		defer cancel()
		f(ctx)
	}()
	return h
}

// Tick starts a task that executes a function at specified intervals.
//
// Parameters:
//   - ctx: The context used to control the lifecycle of the task.
//   - fn:  The function to be executed periodically, accepting a context.
//   - d:   The duration between executions.
//
// Returns:
//   - Handle: A handle that can be used to control the task.
func Tick(ctx context.Context, f func(context.Context), d time.Duration) Handle {
	ctx, cancel := context.WithCancel(ctx)
	h := &taskHandle{
		done:   make(chan struct{}),
		cancel: cancel,
	}

	go func() {
		defer close(h.done)
		defer cancel()
		ticker := time.NewTicker(d)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				f(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
	return h
}
