package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/gojira/ferry/config"
	"github.com/gojira/ferry/jirafinder"
)

var (
	jiraUrl     string
	projectName string
	sprintName  string
	outputFile  string
	configFile  string
)

func init() {
	rootCmd.AddCommand(exportCmd)

	fl := exportCmd.PersistentFlags()

	fl.StringVarP(&configFile, "config", "c", "config.json", "Path to config in json format. default=config.json")
	fl.StringVarP(&outputFile, "output", "o", "", "The target file where output will be exported to")
	fl.StringVar(&jiraUrl, "jira.url", "", "URL to JIRA worskspace, overwrite config.JiraUrl")
	fl.StringVar(&projectName, "project", "", "The project to grab issues from, overwrite config.Filters.Project")
	fl.StringVar(&sprintName, "sprint", "", "Name of the sprint to export, overwrite config.Filters.Sprint")
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Search and export Issues From JIRA",
	RunE: func(cmd *cobra.Command, args []string) error {
		err, c := config.New(configFile)
		if err != nil {
			return err
		}

		//overwrite config
		if outputFile != "" {
			c.DownloadPath = outputFile
		}

		if jiraUrl != "" {
			c.JiraURL = jiraUrl
		}

		if projectName != "" {
			c.Filters["Project"] = projectName
		}

		if sprintName != "" {
			c.Filters["Sprint"] = sprintName
		}

		// start Jira Finder instance
		err, f := jirafinder.NewJiraFinder(c)
		if err != nil {
			return err
		}

		if err := f.Search(); err != nil {
			return err
		}

		fmt.Println(" Download complete!!. Results exported to " + "'" + f.Config.DownloadPath + "'")
		return nil
	},
}
