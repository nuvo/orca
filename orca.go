package main

import (
	"os"

	"github.com/maorfr/orca/cmd"
)

func main() {
	cmd := cmd.NewRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
