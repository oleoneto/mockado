package mocker

import (
	"github.com/cleopatrio/mockado/pkg/envloader"
	"github.com/cleopatrio/mockado/pkg/httpclient"
	"github.com/cleopatrio/mockado/pkg/logger"
	"github.com/cleopatrio/mockado/pkg/types"
)

type mockerFlags struct {
	*types.CLIFlags
	Routes                      []httpclient.Route `yaml:"routes"`
	ShowRoutes                  bool               `yaml:"show_routes"`
	ShowSensitiveVariables      bool               `yaml:"show_sensitive_variables"`
	SensitiveVariablesAllowlist []string           `yaml:"sensitive_variables_allowlist"`
	File                        types.FileFlag     `yaml:"file"`
	Port                        int                `yaml:"port"`
}

type schema struct {
	Routes []httpclient.Route `yaml:"routes"`
}

func init() {
	command.Flags().VarP(&cflags.File, "file", "f", "")
	command.Flags().IntVarP(&cflags.Port, "port", "p", cflags.Port, "")
	command.Flags().BoolVarP(&cflags.ShowRoutes, "show-routes", "r", cflags.ShowRoutes, "")
	command.Flags().BoolVarP(&cflags.ShowRoutes, "show-sensitive-variables", "s", cflags.ShowSensitiveVariables, "")
	// command.Flags().StringArrayVarP(&cflags.SensitiveVariablesAllowlist, "sensitive-variable-allowlist", "l", cflags.SensitiveVariablesAllowlist, "")
}

var cflags = mockerFlags{CLIFlags: types.GetFlags(), Port: 9000}

var log = logger.GetLogger("mocker")

var variableLoader *envloader.Loader
