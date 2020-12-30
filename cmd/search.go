package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gojira/ferry/config"
	"github.com/gojira/ferry/jirafinder"
)

type flagParsingError struct {
	msg string
}

func newFlagParsingError(msg string) *flagParsingError {
	return &flagParsingError{msg: msg}
}
func (e *flagParsingError) Error() string {
	return e.msg
}

var (
	jURL    string
	filters string
	fields  string
	output  string
)

func init() {
	rootCmd.AddCommand(searchCmd)

	fl := searchCmd.PersistentFlags()

	// Required Flags
	fl.StringVarP(&jURL, "url", "u", "", "URL to JIRA worskspace")
	searchCmd.MarkPersistentFlagRequired("url")

	fl.StringVarP(&filters, "filters", "f", "", "The filters to be applied for the search. eg: Project Name, Sprint etc. Filters must be 'quoted' and 'comma' separated with name value 'colon' separated. Eg -> 'Project Name:YourProject,Sprint:Sprint1'")
	searchCmd.MarkPersistentFlagRequired("filters")

	fl.StringVarP(&fields, "fields", "d", "", "The fields to be retrieved for the result set. Fields must be comma separated and quoted")
	searchCmd.MarkPersistentFlagRequired("fields")
}

var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "Search and export Issues From JIRA",
	Example: "ferry search --url yoursite.com --output test.csv --filters 'project name:test project,sprint:sprint 2,issue type:bug' --fields 'assignee,points,scrum team'",
	Long:    "Search and export issues.",
	RunE: func(cmd *cobra.Command, args []string) error {

		// TODO : Remove print for debugging
		fmt.Println(args)
		fmt.Println(filters, fields)

		c, err := gConfig()
		if err != nil {
			return err
		}
		creds, err := gCreds()
		if err != nil {
			return err
		}

		c.SetCreds(creds)

		// start Jira Finder instance
		err, f := jirafinder.NewJiraFinder(c)
		if err != nil {
			return err
		}

		if err := f.Search(); err != nil {
			return err
		}

		fmt.Println(" Download complete!!. Results exported to 'output.csv' file ")
		return nil
	},
}

func gConfig() (*config.Config, error) {

	mFilters, err := parseFilters()
	if err != nil {
		return nil, err
	}
	sFields, err := parseFields()

	if err != nil {
		return nil, err
	}

	return config.NewConfig(jURL, mFilters, sFields, output), nil
}

func gCreds() (config.Creds, error) {

	return config.Creds{Username: "krishnakayc@gmail.com", Password: "3Cw9WbaBp6UD0bLgk0I23AB7"}, nil
	/*  var username string

	fmt.Printf("Please enter credentials for the site %v", jURL)
	fmt.Println("")
	fmt.Print("username: ")
	fmt.Scan(&username)
	fmt.Print("Password (API Token) : ")
	password, err := terminal.ReadPassword(0)

	if err != nil {
		return config.Creds{}, err
	}

	return config.NewCreds(username, string(password)), nil
	*/
}

func parseFilters() (map[string]interface{}, error) {
	mFilters := make(map[string]interface{})
	for _, filter := range strings.Split(filters, ",") {
		kvp := strings.Split(filter, ":")
		if len(kvp) == 2 {
			mFilters[kvp[0]] = kvp[1]
		}
	}

	if len(mFilters) == 0 {
		return nil, newFlagParsingError("Error in 'Filters syntax'. Make sure to enclose --filters flag values with quotes '' and are comma separated. See ferry search --help for more info")
	}

	return mFilters, nil
}

func parseFields() ([]string, error) {
	sFields := make([]string, 0)
	for _, field := range strings.Split(fields, ",") {
		sFields = append(sFields, field)
	}

	if len(sFields) == 0 {
		return nil, newFlagParsingError("Error in 'Fields syntax'. Make sure to enclose --fields flag values with quotes '' and are comma separated. See ferry search --help for more info")
	}
	return sFields, nil

}
