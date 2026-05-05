package main

import (
	"fmt"
	"os"

	"github.com/kasuganosora/bangumi.skill/cli/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	rootCmd := cmd.NewRootCmd(cmd.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
