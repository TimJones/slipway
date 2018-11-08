package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "slip",
	Short: "slipway is a manager od docker for local development",
	Long: `A tool to help make the most of Docker when developing
projects locally`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("slipway")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
