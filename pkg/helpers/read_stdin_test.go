package helpers_test

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/oleoneto/mockado/pkg/helpers"
)

func TestReadStdin(t *testing.T) {
	type test struct {
		name           string
		err            error
		input          io.Reader
		expectedOutput []byte
	}

	tests := []test{
		{
			name:           "test - 1",
			err:            nil,
			input:          bytes.NewReader([]byte(`some-value`)),
			expectedOutput: []byte("some-value"),
		},
		{
			name:           "test - 2",
			err:            fmt.Errorf("invalid reader"),
			input:          bytes.NewReader([]byte("")),
			expectedOutput: []byte(""),
		},
	}

	// FIXME: Need to add some real tests here
	for _, test := range tests {
		output, err := helpers.ReadFromReader(test.input)

		if err != nil {
			if err != test.err {
				t.Errorf("expected %s, but got %s", test.err, err)
			}
		}

		if !reflect.DeepEqual(test.expectedOutput, output) {
			t.Errorf("expected %v, but got %v", string(test.expectedOutput), string(output))
		}
	}
}
