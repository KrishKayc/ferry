package jirafinder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	httprequest "goJIRA/httprequest"
	"os"
	"strings"
)

type configuration struct {
	jiraURL          string
	credentials      credentials
	filters          map[string]interface{}
	fieldsToRetrieve []string
	downloadPath     string
	authToken        string
}

type subTask struct {
	taskType     string
	assigneeName string
	totalHours   string
	name         string
}

type credentials struct {
	username string
	password string
}

type jiraIssue struct {
	data         map[string]interface{}
	subTasks     []subTask
	fields       []string
	assigneeName string
}

// Finder finds the issue from jira based on the config
type finder interface {
	Search()
}

// JiraFinder finds the issue from jira based on the config
type JiraFinder struct {
	config configuration
}

//NewJiraFinder gives a new jira finder with configurations from the config file
func NewJiraFinder(configFile string) JiraFinder {
	return JiraFinder{config: config(configFile)}
}

//Search finds the issue from jira based on the config
func (f JiraFinder) Search() {

	if f.config.jiraURL == "" {
		panic("Set the config first before searching using SetConfig() func")
	}

	output := [][]string{f.config.fieldsToRetrieve}

	out := f.produceFields()

	filters, fields := f.processFields(<-out)

	issues := f.search(filters, fields)

	issue := f.processIssues(<-issues)

	for i := range issue {
		if f := download(i); f != nil {
			output = append(output, f)
		}
	}

	writeToCsv(output, f.config.downloadPath)

	fmt.Println(" Download complete!!. Results exported to " + "'" + f.config.downloadPath + "'")

}

func (f JiraFinder) produceFields() chan []map[string]interface{} {

	out := make(chan []map[string]interface{}, 0)
	go func() {
		req := httprequest.NewHTTPRequest(f.config.jiraURL, "/rest/api/2/field", f.config.authToken, nil)
		body := req.Send()

		var fields []map[string]interface{}
		json.Unmarshal([]byte(body), &fields)

		out <- fields
	}()

	return out

}

func (f JiraFinder) processFields(fields []map[string]interface{}) (map[string]string, []string) {

	filters := make(map[string]string, 0)
	resFields := make([]string, 0)

	var wg sync.WaitGroup
	wg.Add(len(fields))

	for _, field := range fields {
		go func(field map[string]interface{}) {
			defer wg.Done()
			for k, v := range f.config.filters {
				if strings.ToLower(field["name"].(string)) == strings.ToLower(k) {
					key := k
					if field["custom"].(bool) {
						key = "cf[" + strings.Replace(field["id"].(string), "customfield_", "", -1) + "]"
					}
					filters[key] = v.(string)
				}
			}

			for _, v := range f.config.fieldsToRetrieve {
				if strings.ToLower(field["name"].(string)) == strings.ToLower(v) {
					val := v
					if field["custom"].(bool) {
						val = fmt.Sprint(field["id"].(string))
					}
					resFields = append(resFields, val)
				}
			}
		}(field)

	}

	wg.Wait()

	clean(filters)

	return filters, resFields
}

func (f JiraFinder) search(filters map[string]string, fields []string) chan []jiraIssue {
	out := make(chan []jiraIssue, 0)

	go func() {
		params := make(map[string]string, 0)
		params["jql"] = getJql(filters)
		params["fields"] = strings.Join(fields, ",")
		params["maxResults"] = "1000"

		req := httprequest.NewHTTPRequest(f.config.jiraURL, "/rest/api/2/search", f.config.authToken, params)
		body := req.Send()

		var responseResult map[string]interface{}
		var issues []interface{}
		json.Unmarshal([]byte(body), &responseResult)

		issues = responseResult["issues"].([]interface{})

		ji := make([]jiraIssue, 0)
		for _, issue := range issues {
			ji = append(ji, jiraIssue{data: issue.(map[string]interface{}), fields: fields})
		}

		out <- ji
	}()

	return out
}

func (f JiraFinder) processIssues(issues []jiraIssue) chan jiraIssue {

	out := make(chan jiraIssue, 100)
	for i, issue := range issues {
		go func(issue jiraIssue, i int) {
			issueID := issue.data["id"].(string)
			p := f.getIssue(issueID, true)
			parent := <-p
			subTasks := parent["fields"].(map[string]interface{})["subtasks"].([]interface{})
			result := make([]subTask, 0)

			for _, v := range subTasks {
				st := f.getIssue(v.(map[string]interface{})["id"].(string), false)
				subTaskIssue := <-st
				assignee := getValueFromField(subTaskIssue, "assignee")
				issueType := getValueFromField(subTaskIssue, "issuetype")
				name := getValueFromField(subTaskIssue, "summary")
				totalHours := getValueFromField(subTaskIssue, "timetracking")
				currentSubTask := subTask{taskType: issueType, name: name, assigneeName: assignee, totalHours: totalHours}

				result = append(result, currentSubTask)
			}

			issue.subTasks = result

			parentIssueType := getValueFromField(parent, "issuetype")
			if isBug(parentIssueType) {
				issue.assigneeName = getDeveloperNameFromLog(parent)
			}
			out <- issue

			if i == len(issues)-1 {
				close(out)
			}

		}(issue, i)
	}

	return out

}

func (f JiraFinder) getIssue(issueID string, includeChangeLog bool) chan map[string]interface{} {

	out := make(chan map[string]interface{}, 0)
	go func() {
		var getIssueURL string

		if includeChangeLog {
			getIssueURL = "/rest/api/2/issue/" + issueID + "?expand=changelog"
		} else {
			getIssueURL = "/rest/api/2/issue/" + issueID
		}

		req := httprequest.NewHTTPRequest(f.config.jiraURL, getIssueURL, f.config.authToken, nil)
		body := req.Send()

		var responseResult map[string]interface{}
		json.Unmarshal([]byte(body), &responseResult)

		out <- responseResult
	}()

	return out
}

func download(issue jiraIssue) []string {
	fieldValues := make([]string, 0)

	// Listen to final populated issue and prepare the output for all the fields mentioned in the configuration
	for _, field := range issue.fields {
		val, ok := issue.data[field]
		if ok {
			fieldValues = append(fieldValues, strings.Replace(val.(string), ",", "", -1))
		} else {
			fieldValues = append(fieldValues, getFieldValue(field, issue))
		}
	}
	if len(fieldValues) > 0 {
		return fieldValues
	}
	return nil
}

func config(confgFile string) configuration {
	var c map[string]interface{}
	fmt.Println(" Fetching data based on the configuration file => " + "'" + confgFile + "'")
	jsonFile, err := os.Open(confgFile)
	HandleError(err)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &c)
	config := load(c)
	config.authToken = encodeStringToBase64(config.credentials.username + ":" + config.credentials.password)
	return config
}

func load(c map[string]interface{}) configuration {
	config := configuration{jiraURL: c["JiraUrl"].(string), downloadPath: c["DownloadPath"].(string)}
	config.filters = c["Filters"].(map[string]interface{})
	config.fieldsToRetrieve = []string{}
	for _, f := range c["FieldsToRetrieve"].([]interface{}) {
		config.fieldsToRetrieve = append(config.fieldsToRetrieve, f.(string))
	}
	config.credentials = credentials{username: c["Credentials"].(map[string]interface{})["Username"].(string), password: c["Credentials"].(map[string]interface{})["Password"].(string)}
	return config
}
