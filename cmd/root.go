package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var (
	Version = "0.0.1"
)

var rootCmd = &cobra.Command{
	Use:   "ferry",
	Short: "Utility to Search and download Issues From JIRA with a configurable Filters and return fields. Retrieved Issues from JIRA will be downloaded as CSV in the path specified in the config.json file",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
