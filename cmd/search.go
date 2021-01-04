package cmd

import (
	"fmt"
	"strings"

	"github.com/gojira/jira/httprequest"
	"github.com/gojira/jira/search"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
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
	jURL      string
	filters   string
	fields    string
	project   string
	issueType string
)

func init() {
	rootCmd.AddCommand(searchCmd)

	fl := searchCmd.PersistentFlags()

	// Required Flags
	fl.StringVarP(&jURL, "url", "u", "", "URL to JIRA worskspace")
	searchCmd.MarkPersistentFlagRequired("url")

	fl.StringVarP(&project, "project", "p", "",
		"The 'JIRA PROJECT(p)' in which the search has to be done. Give comman separated for multiple projects")
	searchCmd.MarkPersistentFlagRequired("project")

	fl.StringVarP(&issueType, "issuetype", "i", "",
		"The issue type which has to be searched, this can be either a single type or comma separated types eg: Story,Bug")
	searchCmd.MarkPersistentFlagRequired("issueType")

	// Optional Flags
	fl.StringVarP(&filters, "filters", "f", "",
		"The filters to be applied for the search. eg: Project Name, Sprint etc. Filters must be 'quoted' and 'comma' separated with name value 'colon' separated. Eg -> \"Project:YourProject,Sprint:Sprint1\"")

	fl.StringVarP(&fields, "fields", "d", "", "The fields to be retrieved for the result set. Fields must be comma separated and quoted")

}

var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "Search and export Issues From JIRA",
	Example: "jira search --url yoursite.com  --filters \"project:test project,sprint:sprint 2,issue type:bug\" --fields \"assignee,points,scrum team\"",
	Long:    "Search and export issues.",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := gCreds()
		if err != nil {
			return err
		}

		p, err := searchParam(creds)
		if err != nil {
			return err
		}
		if err := search.Search(p, httprequest.NewClient(p.URL, p.AuthToken())); err != nil {
			return err
		}

		fmt.Println("Download complete!!. Results exported to 'output.csv' file ")
		return nil
	},
}

func searchParam(creds search.Creds) (search.Param, error) {

	mFilters, err := parseFilters()
	if err != nil {
		return search.Param{}, err
	}
	//Add project and Issue Type to the filters
	mFilters = append(mFilters, search.NewField("Project", project), search.NewField("Issue Type", issueType))

	sFields, err := parseFields()

	if err != nil {
		return search.Param{}, err
	}

	return search.NewParam(jURL, project, issueType, mFilters, sFields, creds), nil
}

func gCreds() (search.Creds, error) {

	//	return search.Creds{Username: "krishnakayc@gmail.com", Password: "3Cw9WbaBp6UD0bLgk0I23AB7"}, nil
	var username string

	fmt.Printf("Please enter credentials for the site '%v'", jURL)
	fmt.Println("")
	fmt.Print("username: ")
	fmt.Scan(&username)
	fmt.Print("Password (API Token) : ")
	password, err := terminal.ReadPassword(0)
	fmt.Println("")
	if err != nil {
		return search.Creds{}, err
	}

	return search.NewCreds(username, string(password)), nil

}

func parseFilters() ([]search.Field, error) {
	mFilters := make([]search.Field, 0)

	if len(filters) == 0 {
		return mFilters, nil
	}

	for _, filter := range strings.Split(filters, ",") {
		kvp := strings.Split(filter, ":")
		if len(kvp) == 2 {
			mFilters = append(mFilters, search.NewField(kvp[0], kvp[1]))
		}
	}

	if len(mFilters) == 0 {
		return nil, newFlagParsingError("Error in 'Filters syntax'.Make sure to enclose --filters flag values with quotes '' and are comma separated. See ferry search --help for more info")
	}

	return mFilters, nil
}

func parseFields() ([]search.Field, error) {
	sFields := make([]search.Field, 0)
	if len(fields) == 0 {
		//If user hasn't specified any field, add summary and assignee as default fields
		sFields = append(sFields, search.NewField("summary", ""))
		sFields = append(sFields, search.NewField("assignee", ""))

		return sFields, nil
	}
	for _, field := range strings.Split(fields, ",") {
		sFields = append(sFields, search.NewField(field, ""))
	}

	if len(sFields) == 0 {
		return nil, newFlagParsingError("Error in 'Fields syntax'.Make sure to enclose --fields flag values with quotes '' and are comma separated. See ferry search --help for more info")
	}
	return sFields, nil

}
