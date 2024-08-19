package validate

import (
	"errors"
	"slices"

	"github.com/gopherd/core/op"
)

var (
	ErrNotOneOf = errors.New("value is not one of the allowed values")
)

func OneOf[S ~[]T, T comparable](x T, s S) error {
	return op.If(slices.Contains(s, x), nil, ErrNotOneOf)
}
