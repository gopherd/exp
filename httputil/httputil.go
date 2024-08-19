package httputil

import (
	"github.com/gopherd/core/errkit"
)

// Response is a unified response structure for all API endpoints
type Response struct {
	// Error information, if any
	// If this field is not null, it means the request resulted in an error
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message,omitempty"`
	} `json:"error"`

	// The actual data returned by the API
	// This field is populated on successful requests
	Data any `json:"data,omitempty"`
}

// Result returns a Response object from the given value.
func Result(value any) Response {
	if value == nil {
		return Response{}
	}

	if resp, ok := value.(Response); ok {
		return resp
	} else if resp, ok := value.(*Response); ok {
		if resp == nil {
			return Response{}
		}
		return *resp
	}

	if err, ok := value.(error); ok && err != nil {
		var resp Response
		resp.Error.Code = errkit.Errno(err)
		resp.Error.Message = err.Error()
		return resp
	}

	return Response{Data: value}
}

// Binder is an interface for binding request body to data.
type Binder interface {
	// Bind binds the request body to the given data.
	Bind(data any) error
}

// ValueSetter is an interface for setting context value.
type ValueSetter interface {
	// Set sets the value of the given key in the context.
	Set(key string, value any)
}

// Valuer is the interface that wraps the ContextKey method.
type Valuer interface {
	// GetContextKey returns the key of context value without receiver.
	GetContextKey() string
}

// SetContextValue sets the context value to the context.
func SetContextValue[C ValueSetter, V Valuer](ctx C, v V) {
	ctx.Set(v.GetContextKey(), v)
}
