package cmd

import "os"

func Execute() {
	if err := NewRootCommand(os.Stdout, os.Stderr).Execute(); err != nil {
		os.Exit(1)
	}
}
