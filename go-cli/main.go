package main

import (
	"os"

	"gitlab.cc-asp.fraunhofer.de/templates/go-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
