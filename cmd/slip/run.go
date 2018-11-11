package main

import (
	"github.com/spf13/cobra"

	"github.com/timjones/slipway/pkg/slipway"
)

func init() {
	rootCmd.AddCommand(runCommand)
}

var runCommand = &cobra.Command{
	Use:                "run",
	Short:              "Run a command inside the project container",
	DisableFlagParsing: true,
	Run:                run,
}

func run(cmd *cobra.Command, args []string) {
	project := slipway.SlipProject
	if err := project.Build(); err != nil {
		panic(err)
	}
}
