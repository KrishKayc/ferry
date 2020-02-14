package jirafinder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/fatih/color"

	httprequest "jiraSearch_git/httprequest"
	"os"
	"strings"
)

// Configuration represents the configuration file 'config.json'
type Configuration struct {
	JiraURL          string
	Credentials      Credentials
	Filters          map[string]interface{}
	FieldsToRetrieve []string
	DownloadPath     string
	AuthToken        string
}

// SubTask of a Jira Issue
type SubTask struct {
	Type         string
	AssigneeName string
	TotalHours   string
	Name         string
}

// Credentials of the user
type Credentials struct {
	Username string
	Password string
}

// JiraIssue represents an issue in Jira
type JiraIssue struct {
	Data         map[string]interface{}
	SubTasks     []SubTask
	Fields       []string
	AssigneeName string
}

// ProcessedData represents the filters and fields after replacing custom field ids
type ProcessedData struct {
	filters map[string]string
	fields  []string
}

var config Configuration
var workerCount int

// Finder finds the issue from jira based on the config
type Finder interface {
	Search()
}

// JiraFinder finds the issue from jira based on the config
type JiraFinder struct {
}

//Search finds the issue from jira based on the config
func (com *JiraFinder) Search() {

	if config.JiraURL == "" {
		panic("Set the config first before searching using SetConfig() func")
	}
	start := time.Now()

	// Initialize necessary variables
	customFieldChannel := make(chan map[string]string)
	customFieldProcessorChannel := make(chan ProcessedData)
	issueRetrievedChannel := make(chan JiraIssue)
	finalIssueChannel := make(chan JiraIssue)
	outputData := make([][]string, 0)
	issueCount := 0
	totalRestCalls := 0
	totalIssueCount := 0

	// get all the custom fields first and populate in customFieldChannel
	totalRestCalls++
	go getCustomFields(customFieldChannel)

	// listen to custom fields channel and process the data
	go processCustomFields(customFieldChannel, customFieldProcessorChannel)

	// listen to processed data and intiate the search based on it
	go initiateSearch(customFieldProcessorChannel, issueRetrievedChannel, &totalRestCalls)

	//prepare output with the headers first
	headers := config.FieldsToRetrieve
	outputData = append(outputData, headers)

	for {
		select {

		case issue := <-issueRetrievedChannel:
			{
				totalIssueCount++
				issueCount++
				displayProgressAndStatistics(totalIssueCount, -1, totalRestCalls, int(time.Since(start).Seconds()))
				includeChangeLog := strings.Contains(strings.ToLower(config.Filters["IssueType"].(string)), "bug") || strings.Contains(strings.ToLower(config.Filters["IssueType"].(string)), "issue")

				// listen to issue retrieved channel and assign the sub tasks for it
				go getSubTasksForIssue(issue, finalIssueChannel, includeChangeLog, &totalRestCalls)
			}

		case finalIssue := <-finalIssueChannel:
			{
				issueCount--
				fieldValues := make([]string, 0)

				// Listen to final populated issue and prepare the output for all the fields mentioned in the configuration
				for _, field := range finalIssue.Fields {
					val, ok := finalIssue.Data[field]
					if ok {
						fieldValues = append(fieldValues, strings.Replace(val.(string), ",", "", -1))
					} else {
						fieldValues = append(fieldValues, getFieldValue(field, finalIssue))
					}
				}
				if len(fieldValues) > 0 {
					outputData = append(outputData, fieldValues)
				}

				displayProgressAndStatistics(totalIssueCount, issueCount, totalRestCalls, int(time.Since(start).Seconds()))

			}

			if issueCount == 0 {
				// Once all the issues are processed, flush out the output data to the csv file mentioned in the configuration 'DownloadPath'
				writeToCsv(outputData, config.DownloadPath)
				displayProgressAndStatistics(totalIssueCount, issueCount, totalRestCalls, int(time.Since(start).Seconds()))
				fmt.Println()
				fmt.Println()
				color.HiGreen(" Download complete!!. Results exported to " + "'" + config.DownloadPath + "'")
				return
			}

		}
	}
}

// GetCustomFields gets all the custom fields for the jiraUrl mentioned in the config
func getCustomFields(customFieldChannel chan map[string]string) {

	req := httprequest.NewHTTPRequest(config.JiraURL, "/rest/api/2/field", config.AuthToken, nil)
	body := req.Send()

	var fields []map[string]interface{}
	json.Unmarshal([]byte(body), &fields)

	var result map[string]string
	result = make(map[string]string)
	staticFields := make(map[string]string, 0)

	for _, field := range fields {
		if field["custom"].(bool) {
			_, ok := result[field["name"].(string)]
			if !ok {
				_, isStaticField := staticFields[strings.ToLower(field["name"].(string))]
				if !isStaticField {
					result[strings.ToLower(field["name"].(string))] = strings.ToLower(field["id"].(string))
				}
			}
		} else {
			_, ok := result[strings.ToLower(field["name"].(string))]
			if ok {
				delete(result, strings.ToLower(field["name"].(string)))
			} else {
				staticFields[strings.ToLower(field["name"].(string))] = strings.ToLower(field["id"].(string))
			}
		}
	}

	customFieldChannel <- result
}

// ProcessCustomFields replaces custom field names with cf[customFieldId] for searching purpose
// Gets the values from customFieldChannel and replaces the names and assign it to processor channel
func processCustomFields(customFieldChannel chan map[string]string, customFieldProcessorChannel chan ProcessedData) {
	customFieldMap, ok := <-customFieldChannel
	substitutedFilters := make(map[string]string)
	substitutedfieldsToRetrieve := make([]string, 0)

	if ok {
		for k, v := range config.Filters {
			_, ok := customFieldMap[strings.ToLower(k)]
			if ok {
				equivalentFieldID := customFieldMap[strings.ToLower(k)]
				keyName := "cf[" + strings.Replace(equivalentFieldID, "customfield_", "", -1) + "]"
				substitutedFilters[keyName] = v.(string)
			} else {
				substitutedFilters[k] = v.(string)
			}
		}

		for _, v := range config.FieldsToRetrieve {
			_, ok = customFieldMap[strings.ToLower(v)]
			if ok {
				equivalentFieldID := customFieldMap[strings.ToLower(v)]
				substitutedfieldsToRetrieve = append(substitutedfieldsToRetrieve, fmt.Sprint(equivalentFieldID))
			} else {
				substitutedfieldsToRetrieve = append(substitutedfieldsToRetrieve, strings.ToLower(v))
			}
		}
	}
	processedValues := ProcessedData{filters: substitutedFilters, fields: substitutedfieldsToRetrieve}
	customFieldProcessorChannel <- processedValues
}

// InitiateSearch initiates the search process
func initiateSearch(customFieldProcessorChannel chan ProcessedData, issueRetrievedChannel chan JiraIssue, totalRestCalls *int) {
	processedValues, ok := <-customFieldProcessorChannel
	if ok {
		*totalRestCalls++
		go searchIssues(getJql(processedValues.filters), processedValues.fields, issueRetrievedChannel)
	}
}

// GetIssue fetches Issue based from the jiraUrl in the config and issueId passed
func getIssue(issueID string, includeChangeLog bool) map[string]interface{} {

	var getIssueURL string

	if includeChangeLog {
		getIssueURL = "/rest/api/2/issue/" + issueID + "?expand=changelog"
	} else {
		getIssueURL = "/rest/api/2/issue/" + issueID
	}

	req := httprequest.NewHTTPRequest(config.JiraURL, getIssueURL, config.AuthToken, nil)
	body := req.Send()

	var responseResult map[string]interface{}
	json.Unmarshal([]byte(body), &responseResult)

	return responseResult
}

//GetSubTasksForIssue ..
func getSubTasksForIssue(issue JiraIssue, finalIssueChannel chan JiraIssue, includeChangeLog bool, totalRestCalls *int) {

	issueID := issue.Data["id"].(string)
	*totalRestCalls++
	parent := getIssue(issueID, includeChangeLog)
	subTasks := parent["fields"].(map[string]interface{})["subtasks"].([]interface{})
	result := make([]SubTask, 0)

	for _, subTask := range subTasks {
		*totalRestCalls++
		subTaskIssue := getIssue(subTask.(map[string]interface{})["id"].(string), false)
		assignee := getValueFromField(subTaskIssue, "assignee")
		issueType := getValueFromField(subTaskIssue, "issuetype")
		name := getValueFromField(subTaskIssue, "summary")
		totalHours := getValueFromField(subTaskIssue, "timetracking")
		currentSubTask := SubTask{Type: issueType, Name: name, AssigneeName: assignee, TotalHours: totalHours}

		result = append(result, currentSubTask)
	}

	issue.SubTasks = result

	parentIssueType := getValueFromField(parent, "issuetype")
	if isBug(parentIssueType) {
		issue.AssigneeName = getDeveloperNameFromLog(parent)
	}

	finalIssueChannel <- issue

}

func searchIssues(jql string, processedFields []string, issueRetrievedChannel chan JiraIssue) {

	params := make(map[string]string, 0)
	params["jql"] = jql
	params["fields"] = strings.Join(processedFields, ",")
	params["maxResults"] = "1000"

	req := httprequest.NewHTTPRequest(config.JiraURL, "/rest/api/2/search", config.AuthToken, params)
	body := req.Send()

	var responseResult map[string]interface{}
	var issues []interface{}
	json.Unmarshal([]byte(body), &responseResult)

	issues = responseResult["issues"].([]interface{})

	for _, issue := range issues {
		jiraIssue := JiraIssue{Data: issue.(map[string]interface{}), Fields: processedFields}
		issueRetrievedChannel <- jiraIssue
	}

}

//SetConfig ..
func SetConfig(confgFile string) {
	fmt.Println()
	color.Yellow(" Fetching data based on the configuration file => " + "'" + confgFile + "'")
	fmt.Println()
	jsonFile, err := os.Open(confgFile)
	HandleError(err)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &config)

	config.AuthToken = encodeStringToBase64(config.Credentials.Username + ":" + config.Credentials.Password)
	workerCount = 20
}
