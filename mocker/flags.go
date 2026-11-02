package mocker

import (
	"github.com/oleoneto/mockado/pkg/envloader"
	"github.com/oleoneto/mockado/pkg/httpclient"
	"github.com/oleoneto/mockado/pkg/types"
)

type mockerFlags struct {
	*types.CLIFlags
	Routes                 []httpclient.Route `yaml:"routes"`
	ShowRoutes             bool               `yaml:"showRoutes"`
	ShowSensitiveVariables bool               `yaml:"showSensitiveVariables"`
	RedactionPatternsFile  types.FileFlag     `yaml:"sensitiveVariablesBlockList"`
	RoutesFile             types.FileFlag     `yaml:"routesFile"`
	Port                   int                `yaml:"port"`
}

type schema struct {
	Routes []httpclient.Route `yaml:"routes"`
}

func init() {
	command.Flags().VarP(&cflags.RoutesFile, "routes-file", "", "")
	command.Flags().VarP(&cflags.RedactionPatternsFile, "redactor-file", "", "")
	command.Flags().IntVarP(&cflags.Port, "port", "p", cflags.Port, "")
	command.Flags().BoolVarP(&cflags.ShowRoutes, "show-routes", "", cflags.ShowRoutes, "")
	command.Flags().BoolVarP(&cflags.ShowSensitiveVariables, "show-sensitive-variables", "", cflags.ShowSensitiveVariables, "")
}

var cflags = mockerFlags{CLIFlags: types.GetFlags(), Port: 9000}

var loader *envloader.EnvLoader
