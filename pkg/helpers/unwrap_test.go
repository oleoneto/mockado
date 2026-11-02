package helpers_test

import (
	"reflect"
	"testing"

	"github.com/oleoneto/mockado/pkg/helpers"
)

func TestUnwrapOrDefault(t *testing.T) {
	type test[T any] struct {
		name string
		args struct {
			optional *T
			fallback T
		}
		expected T
	}

	t.Run("it should work on built-in types", func(t *testing.T) {
		tests := []test[string]{
			{
				name:     "string - 1",
				expected: "Hello, world!",
				args: struct {
					optional *string
					fallback string
				}{
					optional: nil,
					fallback: "Hello, world!",
				},
			},
			{
				name:     "string - 2",
				expected: "Greetings!",
				args: struct {
					optional *string
					fallback string
				}{
					optional: helpers.PointerTo("Greetings!"),
					fallback: "Hello, world!",
				},
			},
		}

		for _, test := range tests {
			actual := helpers.UnwrapOrDefault(test.args.optional, test.args.fallback)

			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("expected %v but got %v", test.expected, actual)
			}
		}
	})

	t.Run("it should work on structs", func(t *testing.T) {
		type person struct {
			name string
		}

		tests := []test[person]{
			{
				name:     "test - 1",
				expected: person{name: "Unknown"},
				args: struct {
					optional *person
					fallback person
				}{
					optional: nil,
					fallback: person{name: "Unknown"},
				},
			},
			{
				name:     "test - 2",
				expected: person{name: "Leonardo"},
				args: struct {
					optional *person
					fallback person
				}{
					optional: &person{name: "Leonardo"},
					fallback: person{name: "Unknown"},
				},
			},
		}

		for _, test := range tests {
			actual := helpers.UnwrapOrDefault(test.args.optional, test.args.fallback)

			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("expected %v but got %v", test.expected, actual)
			}
		}
	})
}
