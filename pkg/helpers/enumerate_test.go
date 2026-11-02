package helpers_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/oleoneto/mockado/pkg/helpers"
)

func TestEnumerateArgs(t *testing.T) {
	type args struct {
		counter   int
		valueFunc func(index, counter int) string
	}

	type test struct {
		name     string
		args     args
		expected string
	}

	tests := []test{
		{
			name: "test - 1",
			args: args{
				counter: 3,
				valueFunc: func(index int, counter int) string {
					return "-"
				},
			},
			expected: "-, -, -",
		},
		{
			name: "test - 2",
			args: args{
				counter: 3,
				valueFunc: func(index int, counter int) string {
					return fmt.Sprintf(`{%d}`, index)
				},
			},
			expected: "{1}, {2}, {3}",
		},
		{
			name: "test - 3",
			args: args{
				counter: 8,
				valueFunc: func(index int, counter int) string {
					return fmt.Sprintf(`#%d`, index-1)
				},
			},
			expected: "#0, #1, #2, #3, #4, #5, #6, #7",
		},
	}

	for _, test := range tests {
		actual := helpers.EnumerateArgs(test.args.counter, test.args.valueFunc)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("expected %s but got %s", test.expected, actual)
		}
	}
}
