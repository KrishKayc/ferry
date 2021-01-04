package search

import (
	"encoding/base64"
)

//Param represents the search parameters
type Param struct {
	URL       string
	Project   string
	IssueType string
	Creds     Creds
	Filters   []Field
	Fields    []Field
}

//Field is the jira fields which are associated with the issues
type Field struct {
	ID    string
	Name  string
	Value string
}

//User credentials to access the jira
type Creds struct {
	Username string
	Password string
}

//NewCreds creates new credentials instance
func NewCreds(username string, password string) Creds {
	return Creds{Username: username, Password: password}
}

//NewField returns a new field with the name and value
func NewField(name string, value string) Field {
	return Field{Name: name, Value: value}
}

//NewParam returns new search parameters
func NewParam(url string, project string, issueType string, filters []Field, fields []Field, creds Creds) Param {
	return Param{URL: url, Project: project, IssueType: issueType, Filters: filters, Fields: fields, Creds: creds}

}

//AuthToken is used for generating the authentication token for the search parameters
func (s Param) AuthToken() string {
	return encodeStringToBase64(s.Creds.Username + ":" + s.Creds.Password)
}

func encodeStringToBase64(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}
