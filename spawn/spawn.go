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
//   - f:  The task function that accepts a context.
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

type chanOptions struct {
	tickerInterval time.Duration
	tickerFunction func(context.Context)
	cleanup        bool
}

// ChanOption is a configuration option for the Chan functions.
type ChanOption func(*chanOptions)

// WithTicker sets the interval and function for a ticker.
func WithTicker(interval time.Duration, f func(context.Context)) ChanOption {
	if interval <= 0 {
		panic("non-positive interval for WithTicker")
	}
	if f == nil {
		panic("nil function for WithTicker")
	}
	return func(o *chanOptions) {
		o.tickerInterval = interval
		o.tickerFunction = f
	}
}

// WithCleanup specifies whether to clean up the channel after the context is canceled.
func WithCleanup(cleanup bool) ChanOption {
	return func(o *chanOptions) {
		o.cleanup = cleanup
	}
}

func (o *chanOptions) apply(opts []ChanOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func cleanup[T any](ctx context.Context, ch <-chan T, f func(context.Context, T)) {
	for {
		select {
		case v := <-ch:
			f(ctx, v)
		default:
			return
		}
	}
}

// Chan starts a task that processes values from a channel.
func Chan[T any](ctx context.Context, ch <-chan T, f func(context.Context, T), options ...ChanOption) Handle {
	var o chanOptions
	o.apply(options)
	ctx, cancel := context.WithCancel(ctx)
	h := &taskHandle{
		done:   make(chan struct{}),
		cancel: cancel,
	}

	go func() {
		defer close(h.done)
		defer cancel()
		var tc <-chan time.Time
		if o.tickerInterval > 0 {
			ticker := time.NewTicker(o.tickerInterval)
			defer ticker.Stop()
			tc = ticker.C
		}

		for {
			select {
			case <-tc:
				o.tickerFunction(ctx)
			case v := <-ch:
				f(ctx, v)
			case <-ctx.Done():
				if o.cleanup {
					cleanup(ctx, ch, f)
				}
				return
			}
		}
	}()
	return h
}

// Chan2 starts a task that processes values from channel 1 or channel 2.
func Chan2[T1 any, T2 any](ctx context.Context, ch1 <-chan T1, f1 func(context.Context, T1), ch2 <-chan T2, f2 func(context.Context, T2), options ...ChanOption) Handle {
	var o chanOptions
	o.apply(options)
	ctx, cancel := context.WithCancel(ctx)
	h := &taskHandle{
		done:   make(chan struct{}),
		cancel: cancel,
	}

	go func() {
		defer close(h.done)
		defer cancel()
		var tc <-chan time.Time
		if o.tickerInterval > 0 {
			ticker := time.NewTicker(o.tickerInterval)
			defer ticker.Stop()
			tc = ticker.C
		}

		for {
			select {
			case <-tc:
				o.tickerFunction(ctx)
			case v := <-ch1:
				f1(ctx, v)
			case v := <-ch2:
				f2(ctx, v)
			case <-ctx.Done():
				if o.cleanup {
					cleanup(ctx, ch1, f1)
					cleanup(ctx, ch2, f2)
				}
				return
			}
		}
	}()
	return h
}

// Chan3 starts a task that processes values from channel 1, channel 2, or channel 3.
func Chan3[T1 any, T2 any, T3 any](ctx context.Context, ch1 <-chan T1, f1 func(context.Context, T1), ch2 <-chan T2, f2 func(context.Context, T2), ch3 <-chan T3, f3 func(context.Context, T3), options ...ChanOption) Handle {
	var o chanOptions
	o.apply(options)
	ctx, cancel := context.WithCancel(ctx)
	h := &taskHandle{
		done:   make(chan struct{}),
		cancel: cancel,
	}

	go func() {
		defer close(h.done)
		defer cancel()
		var tc <-chan time.Time
		if o.tickerInterval > 0 {
			ticker := time.NewTicker(o.tickerInterval)
			defer ticker.Stop()
			tc = ticker.C
		}

		for {
			select {
			case <-tc:
				o.tickerFunction(ctx)
			case v := <-ch1:
				f1(ctx, v)
			case v := <-ch2:
				f2(ctx, v)
			case v := <-ch3:
				f3(ctx, v)
			case <-ctx.Done():
				if o.cleanup {
					cleanup(ctx, ch1, f1)
					cleanup(ctx, ch2, f2)
					cleanup(ctx, ch3, f3)
				}
				return
			}
		}
	}()
	return h
}

// Chan4 starts a task that processes values from channel 1, channel 2, channel 3, or channel 4.
func Chan4[T1 any, T2 any, T3 any, T4 any](ctx context.Context, ch1 <-chan T1, f1 func(context.Context, T1), ch2 <-chan T2, f2 func(context.Context, T2), ch3 <-chan T3, f3 func(context.Context, T3), ch4 <-chan T4, f4 func(context.Context, T4), options ...ChanOption) Handle {
	var o chanOptions
	o.apply(options)
	ctx, cancel := context.WithCancel(ctx)
	h := &taskHandle{
		done:   make(chan struct{}),
		cancel: cancel,
	}

	go func() {
		defer close(h.done)
		defer cancel()
		var tc <-chan time.Time
		if o.tickerInterval > 0 {
			ticker := time.NewTicker(o.tickerInterval)
			defer ticker.Stop()
			tc = ticker.C
		}

		for {
			select {
			case <-tc:
				o.tickerFunction(ctx)
			case v := <-ch1:
				f1(ctx, v)
			case v := <-ch2:
				f2(ctx, v)
			case v := <-ch3:
				f3(ctx, v)
			case v := <-ch4:
				f4(ctx, v)
			case <-ctx.Done():
				if o.cleanup {
					cleanup(ctx, ch1, f1)
					cleanup(ctx, ch2, f2)
					cleanup(ctx, ch3, f3)
					cleanup(ctx, ch4, f4)
				}
				return
			}
		}
	}()
	return h
}

// Chan5 starts a task that processes values from channel 1, channel 2, channel 3, channel 4, or channel 5.
func Chan5[T1 any, T2 any, T3 any, T4 any, T5 any](ctx context.Context, ch1 <-chan T1, f1 func(context.Context, T1), ch2 <-chan T2, f2 func(context.Context, T2), ch3 <-chan T3, f3 func(context.Context, T3), ch4 <-chan T4, f4 func(context.Context, T4), ch5 <-chan T5, f5 func(context.Context, T5), options ...ChanOption) Handle {
	var o chanOptions
	o.apply(options)
	ctx, cancel := context.WithCancel(ctx)
	h := &taskHandle{
		done:   make(chan struct{}),
		cancel: cancel,
	}

	go func() {
		defer close(h.done)
		defer cancel()
		var tc <-chan time.Time
		if o.tickerInterval > 0 {
			ticker := time.NewTicker(o.tickerInterval)
			defer ticker.Stop()
			tc = ticker.C
		}

		for {
			select {
			case <-tc:
				o.tickerFunction(ctx)
			case v := <-ch1:
				f1(ctx, v)
			case v := <-ch2:
				f2(ctx, v)
			case v := <-ch3:
				f3(ctx, v)
			case v := <-ch4:
				f4(ctx, v)
			case v := <-ch5:
				f5(ctx, v)
			case <-ctx.Done():
				if o.cleanup {
					cleanup(ctx, ch1, f1)
					cleanup(ctx, ch2, f2)
					cleanup(ctx, ch3, f3)
					cleanup(ctx, ch4, f4)
					cleanup(ctx, ch5, f5)
				}
				return
			}
		}
	}()
	return h
}

// Chan6 starts a task that processes values from channel 1, channel 2, channel 3, channel 4, channel 5, or channel 6.
func Chan6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](ctx context.Context, ch1 <-chan T1, f1 func(context.Context, T1), ch2 <-chan T2, f2 func(context.Context, T2), ch3 <-chan T3, f3 func(context.Context, T3), ch4 <-chan T4, f4 func(context.Context, T4), ch5 <-chan T5, f5 func(context.Context, T5), ch6 <-chan T6, f6 func(context.Context, T6), options ...ChanOption) Handle {
	var o chanOptions
	o.apply(options)
	ctx, cancel := context.WithCancel(ctx)
	h := &taskHandle{
		done:   make(chan struct{}),
		cancel: cancel,
	}

	go func() {
		defer close(h.done)
		defer cancel()
		var tc <-chan time.Time
		if o.tickerInterval > 0 {
			ticker := time.NewTicker(o.tickerInterval)
			defer ticker.Stop()
			tc = ticker.C
		}

		for {
			select {
			case <-tc:
				o.tickerFunction(ctx)
			case v := <-ch1:
				f1(ctx, v)
			case v := <-ch2:
				f2(ctx, v)
			case v := <-ch3:
				f3(ctx, v)
			case v := <-ch4:
				f4(ctx, v)
			case v := <-ch5:
				f5(ctx, v)
			case v := <-ch6:
				f6(ctx, v)
			case <-ctx.Done():
				if o.cleanup {
					cleanup(ctx, ch1, f1)
					cleanup(ctx, ch2, f2)
					cleanup(ctx, ch3, f3)
					cleanup(ctx, ch4, f4)
					cleanup(ctx, ch5, f5)
					cleanup(ctx, ch6, f6)
				}
				return
			}
		}
	}()
	return h
}
