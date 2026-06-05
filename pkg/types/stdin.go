package types

import (
	"errors"
	"os"

	"github.com/cleopatrio/mockado/pkg/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type FileFlag struct{ Data []byte }

func (f FileFlag) String() string { return "file" }

func (f FileFlag) Type() string { return "string" }

func (f *FileFlag) Set(value string) (err error) {
	if len(value) < 1 {
		return errors.New("flag not set")
	}

	if f.Data, err = os.ReadFile(value); err != nil {
		return err
	}

	return nil
}

func (f *FileFlag) StdinHook(flagName string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if !cmd.Flag(flagName).Changed {
			data, err := helpers.ReadFromReader(os.Stdin)
			if err != nil {
				os.Exit(1)
			}

			f.Data = data
		}
	}
}

// FileFlag implements pflag.Value
var _ pflag.Value = (*FileFlag)(nil)
