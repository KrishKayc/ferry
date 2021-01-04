package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	Version = "0.0.1"
)

var rootCmd = &cobra.Command{
	Use:   "jira",
	Short: "JIRA searcher. Fetches stories/bugs from your JIRA",
	Long:  "jira CLI to Search and download Issues From JIRA with a configurable Filters and return fields. Retrieved Issues from JIRA will be downloaded as CSV in the path specified in the 'output' flag of the 'search' command.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
}
