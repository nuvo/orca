package main

import (
	"os"

	"orca/cmd"
)

func main() {
	cmd := cmd.NewRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
