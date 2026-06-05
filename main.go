package main

import (
	"os"

	mocker "github.com/cleopatrio/mockado/mocker"
	"github.com/cleopatrio/mockado/pkg/logger"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if err := mocker.Execute(); err != nil {
		os.Exit(1)
	}
}

var log = logger.GetLogger("main")
