package easyecho

import (
	"log/slog"
	"net/http"

	"github.com/gopherd/core/typing"

	"github.com/gopherd/exp/httputil"
)

// Context is an interface for handling HTTP request and response.
type Context interface {
	httputil.Binder
	httputil.ValueSetter

	// JSON sends a JSON response with the given status code and data.
	JSON(statusCode int, resp any) error
	// Get retrieves the value of the given key from the context.
	Get(key string) any
	// Path returns current API path
	Path() string
}

// Router is an interface for registering API endpoints.
type Router[M ~func(H) H, H ~func(C) error, C Context, R any] interface {
	Add(method, path string, handler H, middleware ...M) R
}

// JSON sends a JSON response with the data.
// If the data is nil, it sends a response with empty data.
// If the data is an error, it sends a response with error code and message.
// Otherwise, it sends a response with the data.
func JSON[C Context](ctx C, data any) error {
	return ctx.JSON(http.StatusOK, httputil.Result(data))
}

// BindRequest wraps the handler with request parameter.
func BindRequest[H ~func(C, T) error, C Context, T any](h H) func(C) error {
	return func(ctx C) error {
		var req T
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, typing.Object{"error": err})
			return nil
		}
		return h(ctx, req)
	}
}

// WithValue wraps the handler with context parameter.
func WithValue[H ~func(C, T, V) error, C Context, T any, V httputil.ContextValuer](h H) func(C) error {
	return func(ctx C) error {
		var req T
		if err := ctx.Bind(&req); err != nil {
			slog.Warn("failed to bind request", "error", err, "path", ctx.Path())
			return ctx.JSON(http.StatusBadRequest, typing.Object{"error": err})
		}
		var zero V
		x := ctx.Get(zero.GetContextKey())
		if x == nil {
			slog.Error("context value not found", "path", ctx.Path())
			return ctx.JSON(http.StatusInternalServerError, typing.Object{"error": "context value not found"})
		}
		v, ok := x.(V)
		if !ok {
			slog.Error("unexpected type of context value", "path", ctx.Path())
			return ctx.JSON(http.StatusInternalServerError, typing.Object{"error": "unexpected type of context value"})
		} else {
			return h(ctx, req, v)
		}
	}
}

// Connect adds a CONNECT route to the router.
func Connect[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodConnect, path, BindRequest(f), m...)
}

// Connect2 adds a CONNECT route to the router with context value parameter.
func Connect2[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.ContextValuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodConnect, path, WithValue(f), m...)
}

// Delete adds a DELETE route to the router.
func Delete[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodDelete, path, BindRequest(f), m...)
}

// Delete2 adds a DELETE route to the router with context value parameter.
func Delete2[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.ContextValuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodDelete, path, WithValue(f), m...)
}

// Get adds a GET route to the router.
func Get[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodGet, path, BindRequest(f), m...)
}

// Get2 adds a GET route to the router with context value parameter.
func Get2[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.ContextValuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodGet, path, WithValue(f), m...)
}

// Head adds a HEAD route to the router.
func Head[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodHead, path, BindRequest(f), m...)
}

// Head2 adds a HEAD route to the router with context value parameter.
func Head2[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.ContextValuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodHead, path, WithValue(f), m...)
}

// Options adds a OPTIONS route to the router.
func Options[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodOptions, path, BindRequest(f), m...)
}

// Options2 adds a OPTIONS route to the router with context value parameter.
func Options2[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.ContextValuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodOptions, path, WithValue(f), m...)
}

// Patch adds a PATCH route to the router.
func Patch[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodPatch, path, BindRequest(f), m...)
}

// Patch2 adds a PATCH route to the router with context value parameter.
func Patch2[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.ContextValuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodPatch, path, WithValue(f), m...)
}

// Post adds a POST route to the router.
func Post[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodPost, path, BindRequest(f), m...)
}

// Post2 adds a POST route to the router with context value parameter.
func Post2[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.ContextValuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodPost, path, WithValue(f), m...)
}

// Put adds a PUT route to the router.
func Put[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodPut, path, BindRequest(f), m...)
}

// Put2 adds a PUT route to the router with context value parameter.
func Put2[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.ContextValuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodPut, path, WithValue(f), m...)
}

// Trace adds a TRACE route to the router.
func Trace[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodTrace, path, BindRequest(f), m...)
}

// Trace2 adds a TRACE route to the router with context value parameter.
func Trace2[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.ContextValuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodTrace, path, WithValue(f), m...)
}

// Match adds multiple routes to the router.
func Match[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], methods []string, path string, f F, m ...M) {
	h := BindRequest(f)
	for _, method := range methods {
		router.Add(method, path, h, m...)
	}
}

// Match2 adds multiple routes to the router with context value parameter.
func Match2[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.ContextValuer](router Router[M, H, C, R], methods []string, path string, f F, m ...M) {
	h := WithValue(f)
	for _, method := range methods {
		router.Add(method, path, h, m...)
	}
}
