package search

import (
	"encoding/base64"
)

type SearchParam struct {
	URL       string `json:"JiraUrl"`
	Project   string
	IssueType string
	Creds     Creds   `json:"Credentials"`
	Filters   []Field `json:"Filters"`
	Fields    []Field `json:"FieldsToRetrieve"`
}

type Field struct {
	ID    string
	Name  string
	Value string
}

type Creds struct {
	Username string
	Password string
}

func NewCreds(username string, password string) Creds {
	return Creds{Username: username, Password: password}
}

func NewField(name string, value string) Field {
	return Field{Name: name, Value: value}
}
func NewSearchParam(url string, project string, issueType string, filters []Field, fields []Field, creds Creds) SearchParam {
	return SearchParam{URL: url, Project: project, IssueType: issueType, Filters: filters, Fields: fields, Creds: creds}

}
func (s SearchParam) AuthToken() string {
	return encodeStringToBase64(s.Creds.Username + ":" + s.Creds.Password)
}

func encodeStringToBase64(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}
