package mocker

import (
	"bytes"
	"fmt"

	"github.com/oleoneto/mockado/pkg/envloader"
	"github.com/oleoneto/mockado/pkg/httpclient"
	"github.com/oleoneto/mockado/pkg/logger"
	"github.com/oleoneto/mockado/pkg/redactor"
	"github.com/oleoneto/mockado/pkg/server"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

var logr = logger.GetLogger()

var command = &cobra.Command{
	Use:   "mocker",
	Short: "Run a mock webserver based on a provided file",
	PreRun: func(cmd *cobra.Command, args []string) {
		cflags.RoutesFile.StdinHook("routes-file")
		cflags.RoutesFile.StdinHook("redactor-file")
	},
	Run: func(cmd *cobra.Command, args []string) {
		data, err := loadRoutes()
		if err != nil {
			panic(err)
		}

		var patterns []redactor.SecretPattern
		if !cflags.ShowSensitiveVariables {
			patterns, err = loadRedactorPatterns()
			if err != nil {
				panic(err)
			}
		}

		r := redactor.New(
			redactor.WithLogger(logr),
			redactor.WithRedactedPattern("[redacted]"),
			redactor.WithBlockedPatterns(patterns...),
		)

		l := envloader.New(
			envloader.WithLogger(logr),
			envloader.WithRedactor(r),
		)

		server := server.New(
			server.WithEnvLoader(l),
			server.WithRoutes(data.Routes...),
		)

		server.Listen(cflags.Port, cflags.ShowRoutes)
	},
}

func Execute() error { return command.Execute() }

func loadRoutes() (data schema, err error) {
	bf := bytes.NewBuffer(cflags.RoutesFile.Data)
	loader, err := yaml.NewLoader(bf)
	if err != nil {
		return
	}

	type routesFile struct {
		Routes []httpclient.Route `yaml:"routes" json:"routes"`
	}

	for {
		var file routesFile

		err := loader.Load(&file)
		if err != nil {
			break
		}

		data.Routes = append(data.Routes, file.Routes...)
	}

	if len(data.Routes) == 0 {
		err = fmt.Errorf("No routes found in the provided file")
	}

	return data, nil
}

func loadRedactorPatterns() (data []redactor.SecretPattern, err error) {
	bf := bytes.NewBuffer(cflags.RedactionPatternsFile.Data)
	loader, err := yaml.NewLoader(bf)
	if err != nil {
		return
	}

	type patternsFile struct {
		Patterns []redactor.SecretPattern `yaml:"patterns" json:"patterns"`
	}

	for {
		var file patternsFile

		err := loader.Load(&file)
		if err != nil {
			// logr.Errorln(err)
			break
		}

		data = append(data, file.Patterns...)
	}

	if len(data) == 0 {
		err = fmt.Errorf("No patterns loaded from file")
	}

	return data, err
}
