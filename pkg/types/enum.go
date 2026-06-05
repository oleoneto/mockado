package types

import (
	"fmt"
	"slices"

	"github.com/spf13/pflag"
)

type FlagEnum[T comparable] struct {
	Allowed []T
	Default T
	value   T
}

func (e FlagEnum[T]) Type() string { return "" }

func (e FlagEnum[string]) String() string { return e.value }

func (e *FlagEnum[T]) Set(value T) error {
	if !slices.Contains(e.Allowed, value) {
		return fmt.Errorf("%v is not a supported value", value)
	}

	e.value = value
	return nil
}

// FlagEnum implements pflag.Value
var _ pflag.Value = (*FlagEnum[string])(nil)
