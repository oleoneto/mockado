package mocker

import (
	"bytes"

	"github.com/cleopatrio/mockado/pkg/envloader"
	"github.com/cleopatrio/mockado/pkg/httpclient"
	"github.com/cleopatrio/mockado/pkg/redactor"
	"github.com/cleopatrio/mockado/pkg/server"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

var command = &cobra.Command{
	Use:    "mocker",
	Short:  "Run a mock webserver based on a provided file",
	PreRun: cflags.File.StdinHook("file"),
	Run: func(cmd *cobra.Command, args []string) {
		data, _ := load()

		variableLoader = envloader.WithRedactor(redactor.New())

		server := server.New()

		server.Add(data.Routes...)

		server.Listen(cflags.Port)
	},
}

func Execute() error { return command.Execute() }

func load() (schema, error) {
	var data schema

	bf := bytes.NewBuffer(cflags.File.Data)
	loader, err := yaml.NewLoader(bf)
	if err != nil {
		log.Fatalln("Unmarshall error:", err)
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
		log.Fatalln("No routes found in the provided file")
	}

	return data, nil
}
