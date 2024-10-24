package easygin

import (
	"log/slog"
	"net/http"

	"github.com/gopherd/core/types"

	"github.com/gopherd/exp/httputil"
)

// Context is an interface for handling HTTP request and response.
type Context interface {
	httputil.Binder
	httputil.ValueSetter

	// JSON sends a JSON response with the given status code and data.
	JSON(statusCode int, resp any)
	// Get retrieves the value of the given key from the context.
	Get(key string) (any, bool)
	// FullPath returns current API path
	FullPath() string
}

// Router is an interface for registering API endpoints.
type Router[H ~func(C), C Context, R any] interface {
	Handle(method, path string, handlers ...H) R
}

// JSON sends a JSON response with the data.
// If the data is nil, it sends a response with empty data.
// If the data is an error, it sends a response with error code and message.
// Otherwise, it sends a response with the data.
func JSON[C Context](ctx C, data any) {
	ctx.JSON(http.StatusOK, httputil.Result(data))
}

// BindRequest wraps the handler with request parameter.
func BindRequest[H ~func(C, T), C Context, T any](h H) func(C) {
	return func(ctx C) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, types.Object{"error": err})
			return
		}
		h(ctx, req)
	}
}

// WithValue wraps the handler with context parameter.
func WithValue[H ~func(C, T, V), C Context, T any, V httputil.ContextValuer](h H) func(C) {
	return func(ctx C) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			slog.Warn("failed to bind request", "error", err, "path", ctx.FullPath())
			ctx.JSON(http.StatusBadRequest, types.Object{"error": err})
			return
		}
		var zero V
		x, ok := ctx.Get(zero.GetContextKey())
		if !ok {
			slog.Error("context value not found", "path", ctx.FullPath())
			ctx.JSON(http.StatusInternalServerError, types.Object{"error": "context value not found"})
			return
		}
		if v, ok := x.(V); !ok {
			slog.Error("unexpected type of context value", "path", ctx.FullPath())
			ctx.JSON(http.StatusInternalServerError, types.Object{"error": "unexpected type of context value"})
		} else {
			h(ctx, req, v)
		}
	}
}

// Connect adds a CONNECT route to the router.
func Connect[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodConnect, path, BindRequest(f))
}

// Connect2 adds a CONNECT route to the router with context value parameter.
func Connect2[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.ContextValuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodConnect, path, WithValue(f))
}

// Delete adds a DELETE route to the router.
func Delete[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodDelete, path, BindRequest(f))
}

// Delete2 adds a DELETE route to the router with context value parameter.
func Delete2[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.ContextValuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodDelete, path, WithValue(f))
}

// Get adds a GET route to the router.
func Get[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodGet, path, BindRequest(f))
}

// Get2 adds a GET route to the router with context value parameter.
func Get2[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.ContextValuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodGet, path, WithValue(f))
}

// Head adds a HEAD route to the router.
func Head[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodHead, path, BindRequest(f))
}

// Head2 adds a HEAD route to the router with context value parameter.
func Head2[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.ContextValuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodHead, path, WithValue(f))
}

// Options adds a OPTIONS route to the router.
func Options[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodOptions, path, BindRequest(f))
}

// Options2 adds a OPTIONS route to the router with context value parameter.
func Options2[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.ContextValuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodOptions, path, WithValue(f))
}

// Patch adds a PATCH route to the router.
func Patch[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodPatch, path, BindRequest(f))
}

// Patch2 adds a PATCH route to the router with context value parameter.
func Patch2[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.ContextValuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodPatch, path, WithValue(f))
}

// Post adds a POST route to the router.
func Post[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodPost, path, BindRequest(f))
}

// Post2 adds a POST route to the router with context value parameter.
func Post2[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.ContextValuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodPost, path, WithValue(f))
}

// Put adds a PUT route to the router.
func Put[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodPut, path, BindRequest(f))
}

// Put2 adds a PUT route to the router with context value parameter.
func Put2[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.ContextValuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodPut, path, WithValue(f))
}

// Trace adds a TRACE route to the router.
func Trace[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodTrace, path, BindRequest(f))
}

// Trace2 adds a TRACE route to the router with context value parameter.
func Trace2[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.ContextValuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodTrace, path, WithValue(f))
}

// Match adds multiple routes to the router.
func Match[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], methods []string, path string, f F) {
	h := BindRequest(f)
	for _, method := range methods {
		router.Handle(method, path, h)
	}
}

// Match2 adds multiple routes to the router with context value parameter.
func Match2[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.ContextValuer](router Router[H, C, R], methods []string, path string, f F) {
	h := WithValue(f)
	for _, method := range methods {
		router.Handle(method, path, h)
	}
}
