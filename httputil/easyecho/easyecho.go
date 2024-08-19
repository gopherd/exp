package easyecho

import (
	"net/http"

	"github.com/gopherd/core/types"

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
}

// Router is an interface for registering API endpoints.
type Router[M ~func(H) H, H ~func(C) error, C Context, R any] interface {
	Add(method, path string, handler H, middleware ...M) R
}

// JSON sends a JSON response with the result.
// If the result is nil, it sends a response with empty data.
// If the result is an error, it sends a response with error code and message.
// Otherwise, it sends a response with the result.
func JSON[C Context](ctx C, data any) error {
	return ctx.JSON(http.StatusOK, httputil.Result(data))
}

// BindRequest wraps the handler with request parameter.
func BindRequest[H ~func(C, T) error, C Context, T any](h H) func(C) error {
	return func(ctx C) error {
		var req T
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, types.Object{"error": err})
			return nil
		}
		return h(ctx, req)
	}
}

// WithValue wraps the handler with context parameter.
func WithValue[H ~func(C, T, V) error, C Context, T any, V httputil.Valuer](h H) func(C) error {
	return func(ctx C) error {
		var req T
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, types.Object{"error": err})
			return nil
		}
		var zero V
		v, ok := ctx.Get(zero.GetContextKey()).(V)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, types.Object{"error": "context value not found"})
			return nil
		}
		return h(ctx, req, v)
	}
}

// Connect adds a CONNECT route to the router.
func Connect[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodConnect, path, BindRequest(f), m...)
}

// ConnectWithValue adds a CONNECT route to the router with context value parameter.
func ConnectWithValue[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.Valuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodConnect, path, WithValue(f), m...)
}

// Delete adds a DELETE route to the router.
func Delete[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodDelete, path, BindRequest(f), m...)
}

// DeleteWithValue adds a DELETE route to the router with context value parameter.
func DeleteWithValue[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.Valuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodDelete, path, WithValue(f), m...)
}

// Get adds a GET route to the router.
func Get[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodGet, path, BindRequest(f), m...)
}

// GetWithValue adds a GET route to the router with context value parameter.
func GetWithValue[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.Valuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodGet, path, WithValue(f), m...)
}

// Head adds a HEAD route to the router.
func Head[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodHead, path, BindRequest(f), m...)
}

// HeadWithValue adds a HEAD route to the router with context value parameter.
func HeadWithValue[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.Valuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodHead, path, WithValue(f), m...)
}

// Options adds a OPTIONS route to the router.
func Options[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodOptions, path, BindRequest(f), m...)
}

// OptionsWithValue adds a OPTIONS route to the router with context value parameter.
func OptionsWithValue[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.Valuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodOptions, path, WithValue(f), m...)
}

// Patch adds a PATCH route to the router.
func Patch[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodPatch, path, BindRequest(f), m...)
}

// PatchWithValue adds a PATCH route to the router with context value parameter.
func PatchWithValue[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.Valuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodPatch, path, WithValue(f), m...)
}

// Post adds a POST route to the router.
func Post[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodPost, path, BindRequest(f), m...)
}

// PostWithValue adds a POST route to the router with context value parameter.
func PostWithValue[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.Valuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodPost, path, WithValue(f), m...)
}

// Put adds a PUT route to the router.
func Put[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodPut, path, BindRequest(f), m...)
}

// PutWithValue adds a PUT route to the router with context value parameter.
func PutWithValue[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.Valuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodPut, path, WithValue(f), m...)
}

// Trace adds a TRACE route to the router.
func Trace[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodTrace, path, BindRequest(f), m...)
}

// TraceWithValue adds a TRACE route to the router with context value parameter.
func TraceWithValue[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.Valuer](router Router[M, H, C, R], path string, f F, m ...M) {
	router.Add(http.MethodTrace, path, WithValue(f), m...)
}

// Match adds multiple routes to the router.
func Match[F func(C, T) error, M ~func(H) H, H ~func(C) error, C Context, R, T any](router Router[M, H, C, R], methods []string, path string, f F, m ...M) {
	h := BindRequest(f)
	for _, method := range methods {
		router.Add(method, path, h, m...)
	}
}

// MatchWithValue adds multiple routes to the router with context value parameter.
func MatchWithValue[F func(C, T, V) error, M ~func(H) H, H ~func(C) error, C Context, R, T any, V httputil.Valuer](router Router[M, H, C, R], methods []string, path string, f F, m ...M) {
	h := WithValue(f)
	for _, method := range methods {
		router.Add(method, path, h, m...)
	}
}
