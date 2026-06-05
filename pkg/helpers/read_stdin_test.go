package helpers_test

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/cleopatrio/mockado/pkg/helpers"
)

func TestReadStdin(t *testing.T) {
	type test struct {
		name   string
		err    error
		input  io.Reader
		output []byte
	}

	tests := []test{
		{
			name:   "test - 1",
			err:    nil,
			input:  bytes.NewReader([]byte(`some-value`)),
			output: []byte("some-value"),
		},
		{
			name:   "test - 2",
			err:    fmt.Errorf("invalid reader"),
			input:  bytes.NewReader([]byte("{null}")),
			output: []byte("{null}"),
		},
	}

	for _, test := range tests {
		output, err := helpers.ReadFromReader(test.input)

		if (err != nil) != (test.err != nil) {
			t.Errorf("expected %s, but got %s", test.err, err)
		}

		if !reflect.DeepEqual(test.output, output) {
			t.Errorf("expected %v, but got %v", string(test.output), string(output))
		}
	}
}
