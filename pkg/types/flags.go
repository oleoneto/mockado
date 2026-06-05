package types

import (
	"embed"
	"os"
	"reflect"
	"sync"
	"time"

	gJSON "github.com/drewstinnett/gout/v2/formats/json"
	gYAML "github.com/drewstinnett/gout/v2/formats/yaml"

	"github.com/drewstinnett/gout/v2"
	"github.com/drewstinnett/gout/v2/config"
	"github.com/drewstinnett/gout/v2/formats/gotemplate"

	"github.com/spf13/cobra"
)

var flagsOnce sync.Once

var outputFormat = &FlagEnum[string]{
	Allowed: []string{"json", "yaml", "table", "silent"},
	Default: "yaml",
}

var globalFlags *CLIFlags

type CLIFlags struct {
	BaseDir          string            `yaml:"base_directory"`  // CLI's home directory inside the current user's directory
	ConfigFile       string            `yaml:"config_file"`     // Config file
	IsVerboseEnabled bool              `yaml:"verbose.enabled"` // Enables detailed logging of commands
	IsTimingEnabled  bool              `yaml:"timing.enabled"`  // Calculates and exhibits how long commands take to finish executing
	OutputFormat     *FlagEnum[string] `yaml:"output.format"`   // Encoding format for command output
	OutputTemplate   string            `yaml:"output.template"`

	Plugins            map[string]reflect.Value // External modules
	ExecutionStartTime time.Time
	ExecutionExitLog   []any

	FileSystem embed.FS

	Writer *gout.Gout
}

func GetFlags() *CLIFlags {
	flagsOnce.Do(func() {
		if globalFlags != nil {
			return
		}

		home, _ := os.UserHomeDir()
		globalFlags = &CLIFlags{
			Writer:       gout.New(),
			OutputFormat: outputFormat,
			BaseDir:      home + "/.httpy",
		}
	})

	return globalFlags
}

func (state *CLIFlags) SetFormatter(cmd *cobra.Command, args []string) {
	outputFormat := cmd.Flag("output").Value.String()
	switch outputFormat {
	case "table":
		state.Writer.SetFormatter(&TableFormatter{})
	case "json":
		state.Writer.SetFormatter(gJSON.Formatter{})
	case "yaml":
		state.Writer.SetFormatter(gYAML.Formatter{})
	case "gotemplate":
		state.Writer.SetFormatter(gotemplate.Formatter{
			Opts: config.FormatterOpts{"template": state.OutputTemplate},
		})
	case "silent":
		state.Writer.SetFormatter(SilentFormatter{})
	case "plain":
		state.Writer.SetFormatter(PlainFormatter{})
	}
}
