package cmd

import (
	"fmt"
	"log"

	"github.com/gojira/jira/search"
	"github.com/spf13/cobra"
)

var writer search.Writer
var rootCmd = &cobra.Command{
	Use:   "jira",
	Short: "JIRA searcher. Fetches stories/bugs from your JIRA",
	Long:  "jira CLI to Search and download Issues From JIRA with a configurable Filters and return fields. Retrieved Issues from JIRA will be downloaded as CSV in the path specified in the 'output' flag of the 'search' command.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

//Execute runs the root command
func Execute(w search.Writer) {
	writer = w
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
}
