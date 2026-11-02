package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	mocker "github.com/oleoneto/mockado/mocker"
)

func main() {
	if err := mocker.Execute(); err != nil {
		os.Exit(1)
	}
}
