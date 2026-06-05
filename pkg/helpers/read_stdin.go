package helpers

import (
	"io"
	"strings"
)

func ReadFromReader(f io.Reader) ([]byte, error) { return io.ReadAll(f) }

func ReadStringFromReader(f io.Reader) (string, error) {
	stdin, err := ReadFromReader(f)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdin)), nil
}
