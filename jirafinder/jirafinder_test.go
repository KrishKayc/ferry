package jirafinder

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJiraFinder_CreateConfigEmpty(t *testing.T) {
	r := assert.New(t)

	err, _ := createConfig("")
	r.Errorf(err, "expected createConfig to fail")
	r.Containsf(err.Error(), "empty config not allowed", "expected 'empty config not allowed', got '%s'", err)
}

func TestJiraFinder_CreateConfigNotExisting(t *testing.T) {
	r := assert.New(t)

	err, _ := createConfig("/nofile/here")
	r.Errorf(err, "expected createConfig to fail")
	r.Containsf(err.Error(), "no such file", "expected 'Not Exists Error', got '%s'", err)
}

func TestJiraFinder_CreateConfigSuccess(t *testing.T) {
	r := assert.New(t)
	err, config := createConfig("../example_config/sample_config_bug_search.json")

	r.NoErrorf(err, "expected reading config succeed, got error: '%s'", err)
	r.NotNil(config, "expected to have an healthy config, got nil")
	r.NotEmpty(config.AuthToken, ".AuthToken should not be empty")
}

func TestJiraFinder_DownloadIssue(t *testing.T) {
	r := assert.New(t)

	issue := JiraIssue{
		Data: map[string]interface{}{
			"key":    "POS-7",
			"id":     "10006",
			"expand": "operations,versionedRepresentations,editmeta,changelog,renderedFields",
			"fields": map[string]interface{}{
				"summary":           "Fix issue",
				"customfield_10026": "null",
			},
		},
		Fields: []string{"key", "summary", "assignee"},
	}

	row := download(issue)
	expectedValue := []string{"POS-7", "Fix issue", "N/A"}
	r.EqualValues(expectedValue, row, "Wrong result")
}

func TestJiraFinder_DownloadIssueEmpty(t *testing.T) {
	r := assert.New(t)
	issue := JiraIssue{
		Data: map[string]interface{}{
			"key":    "POS-7",
			"id":     "10006",
		},
		Fields: []string{},
	}

	row := download(issue)
	r.EqualValues([]string{}, row, "Expected empty row")
}

func TestJiraFinder_NewFinder(t *testing.T) {
	r := require.New(t)

	var err error
	var f *JiraFinder

	err, f = NewJiraFinder("../example_config/sample_for_test.json")
	r.NoErrorf(err, "instantiation resulting to error: '%s'", err)
	r.NotNil(f, "finder object nil")
	r.EqualValues("https://your-jira-url.com", f.Config.JiraURL, "wrong jira endpoint")
}

func TestJiraFinder_Search(t *testing.T) {
	r := require.New(t)
	err, f := NewJiraFinder("../example_config/sample_for_test.json")
	r.NoErrorf(err, "instantiation resulting to error: '%s'", err)
	r.NotNil(f, "finder object nil")

	f.UseStub()

	err = f.Search()
	r.NoErrorf(err, "search func resulting to error: %s", err)
}
