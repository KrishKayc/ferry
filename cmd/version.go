package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ferry - JIRA search\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("GoVersion: %s\n", runtime.Version())
	},
}
