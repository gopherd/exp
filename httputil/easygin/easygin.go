package easygin

import (
	"net/http"

	"github.com/gopherd/core/types"

	"reinforce.life/server/pkg/httputil"
)

// Context is an interface for handling HTTP request and response.
type Context interface {
	httputil.Binder
	httputil.ValueSetter

	// JSON sends a JSON response with the given status code and data.
	JSON(statusCode int, resp any)
	// Get retrieves the value of the given key from the context.
	Get(key string) (any, bool)
}

// Router is an interface for registering API endpoints.
type Router[H ~func(C), C Context, R any] interface {
	Handle(method, path string, handlers ...H) R
}

// JSON sends a JSON response with the result.
// If the result is nil, it sends a response with empty data.
// If the result is an error, it sends a response with error code and message.
// Otherwise, it sends a response with the result.
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
func WithValue[H ~func(C, T, V), C Context, T any, V httputil.Valuer](h H) func(C) {
	return func(ctx C) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, types.Object{"error": err})
			return
		}
		var zero V
		v, ok := ctx.Get(zero.GetContextKey())
		if !ok {
			ctx.JSON(http.StatusInternalServerError, types.Object{"error": "context value not found"})
			return
		}
		h(ctx, req, v.(V))
	}
}

// Connect adds a CONNECT route to the router.
func Connect[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodConnect, path, BindRequest(f))
}

// ConnectWithValue adds a CONNECT route to the router with context value parameter.
func ConnectWithValue[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.Valuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodConnect, path, WithValue(f))
}

// Delete adds a DELETE route to the router.
func Delete[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodDelete, path, BindRequest(f))
}

// DeleteWithValue adds a DELETE route to the router with context value parameter.
func DeleteWithValue[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.Valuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodDelete, path, WithValue(f))
}

// Get adds a GET route to the router.
func Get[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodGet, path, BindRequest(f))
}

// GetWithValue adds a GET route to the router with context value parameter.
func GetWithValue[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.Valuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodGet, path, WithValue(f))
}

// Head adds a HEAD route to the router.
func Head[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodHead, path, BindRequest(f))
}

// HeadWithValue adds a HEAD route to the router with context value parameter.
func HeadWithValue[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.Valuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodHead, path, WithValue(f))
}

// Options adds a OPTIONS route to the router.
func Options[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodOptions, path, BindRequest(f))
}

// OptionsWithValue adds a OPTIONS route to the router with context value parameter.
func OptionsWithValue[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.Valuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodOptions, path, WithValue(f))
}

// Patch adds a PATCH route to the router.
func Patch[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodPatch, path, BindRequest(f))
}

// PatchWithValue adds a PATCH route to the router with context value parameter.
func PatchWithValue[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.Valuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodPatch, path, WithValue(f))
}

// Post adds a POST route to the router.
func Post[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodPost, path, BindRequest(f))
}

// PostWithValue adds a POST route to the router with context value parameter.
func PostWithValue[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.Valuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodPost, path, WithValue(f))
}

// Put adds a PUT route to the router.
func Put[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodPut, path, BindRequest(f))
}

// PutWithValue adds a PUT route to the router with context value parameter.
func PutWithValue[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.Valuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodPut, path, WithValue(f))
}

// Trace adds a TRACE route to the router.
func Trace[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodTrace, path, BindRequest(f))
}

// TraceWithValue adds a TRACE route to the router with context value parameter.
func TraceWithValue[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.Valuer](router Router[H, C, R], path string, f F) {
	router.Handle(http.MethodTrace, path, WithValue(f))
}

// Match adds multiple routes to the router.
func Match[F func(C, T), H ~func(C), C Context, R, T any](router Router[H, C, R], methods []string, path string, f F) {
	h := BindRequest(f)
	for _, method := range methods {
		router.Handle(method, path, h)
	}
}

// MatchWithValue adds multiple routes to the router with context value parameter.
func MatchWithValue[F func(C, T, V), H ~func(C), C Context, R, T any, V httputil.Valuer](router Router[H, C, R], methods []string, path string, f F) {
	h := WithValue(f)
	for _, method := range methods {
		router.Handle(method, path, h)
	}
}
