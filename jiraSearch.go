package main

import (
	"fmt"
	"strings"
	"time"

	ui "github.com/gizak/termui"
)

// ProcessedData represents the filters and fields after replacing custom field ids
type ProcessedData struct {
	filters map[string]string
	fields  []string
}

// Configuration represents the configuration file 'config.json'
type Configuration struct {
	JiraUrl          string
	Credentials      Credentials
	Filters          map[string]interface{}
	FieldsToRetrieve []string
	DownloadPath     string
	AuthToken        string
}

// Credentials of the user
type Credentials struct {
	Username string
	Password string
}

// SubTask of a Jira Issue
type SubTask struct {
	Type         string
	AssigneeName string
	TotalHours   string
	Name         string
}

// JiraIssue represents an issue in Jira
type JiraIssue struct {
	Data         map[string]interface{}
	SubTasks     []SubTask
	Fields       []string
	AssigneeName string
}

func main() {
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

	// Init the termui
	// for displaying the output statistics in terminal
	err := ui.Init()
	HandleError(err)
	defer ui.Close()
	g, bc := GetProgressAndStatisticsBar()

	// Read config from the 'config.json' file
	config := ReadConfig("config.json")
	communicator := JiraCommunicator{Url: config.JiraUrl, AuthToken: config.AuthToken}

	// get all the custom fields first and populate in customFieldChannel
	totalRestCalls++
	go GetCustomFields(config, customFieldChannel, &communicator)

	// listen to custom fields channel and process the data
	go ProcessCustomFields(config, customFieldChannel, customFieldProcessorChannel)

	// listen to processed data and intiate the search based on it
	go InitiateSearch(config, customFieldProcessorChannel, issueRetrievedChannel, &totalRestCalls, &communicator)

	//prepare output with the headers first
	headers := config.FieldsToRetrieve
	outputData = append(outputData, headers)

	for {
		select {

		case issue := <-issueRetrievedChannel:
			{
				totalIssueCount++
				issueCount++
				DisplayProgressAndStatistics(totalIssueCount, -1, totalRestCalls, int(time.Since(start).Seconds()), g, bc)
				includeChangeLog := strings.Contains(strings.ToLower(config.Filters["IssueType"].(string)), "bug") || strings.Contains(strings.ToLower(config.Filters["IssueType"].(string)), "issue")

				// listen to issue retrieved channel and assign the sub tasks for it
				go GetSubTasksForIssue(config, issue, finalIssueChannel, includeChangeLog, &totalRestCalls, &communicator)
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
						fieldValues = append(fieldValues, GetFieldValue(field, finalIssue))
					}
				}
				if len(fieldValues) > 0 {
					outputData = append(outputData, fieldValues)
				}

				DisplayProgressAndStatistics(totalIssueCount, issueCount, totalRestCalls, int(time.Since(start).Seconds()), g, bc)

			}

			if issueCount == 0 {
				// Once all the issues are processed, flush out the output data to the csv file mentioned in the configuration 'DownloadPath'
				WriteToCsv(outputData, config.DownloadPath)
				DisplayProgressAndStatistics(totalIssueCount, issueCount, totalRestCalls, int(time.Since(start).Seconds()), g, bc)
				fmt.Println("Download complete!!")
				break
			}

		}
	}
}

// GetJql constructs the JQL string based on filters in the configuration
func GetJql(filters map[string]string) string {
	index := 0
	totalCount := len(filters)
	var b strings.Builder
	for k, v := range filters {
		index++
		if strings.Contains(v, ",") {
			valSlice := strings.Split(v, ",")
			b.WriteString(k + " in (" + GetInFilterValue(valSlice) + ")")
		} else {
			b.WriteString(k + "=" + "'" + v + "'")
		}

		if index != totalCount {
			b.WriteString(" AND ")
		}

	}

	return b.String()
}

// GetInFilterValue constructs the 'In' Clause value
func GetInFilterValue(values []string) string {
	index := 0
	totalCount := len(values)
	var b strings.Builder
	for _, val := range values {
		index++
		b.WriteString("'" + strings.TrimSpace(val) + "'")
		if index != totalCount {
			b.WriteString(",")
		}
	}

	return b.String()
}

// ProcessCustomFields replaces custom field names with cf[customFieldId] for searching purpose
// Gets the values from customFieldChannel and replaces the names and assign it to processor channel
func ProcessCustomFields(config Configuration, customFieldChannel chan map[string]string, customFieldProcessorChannel chan ProcessedData) {
	customFieldMap, ok := <-customFieldChannel
	substitutedFilters := make(map[string]string)
	substitutedfieldsToRetrieve := make([]string, 0)

	if ok {
		for k, v := range config.Filters {
			_, ok := customFieldMap[strings.ToLower(k)]
			if ok {
				equivalentFieldId := customFieldMap[strings.ToLower(k)]
				keyName := "cf[" + strings.Replace(equivalentFieldId, "customfield_", "", -1) + "]"
				substitutedFilters[keyName] = v.(string)
			} else {
				substitutedFilters[k] = v.(string)
			}
		}

		for _, v := range config.FieldsToRetrieve {
			_, ok = customFieldMap[strings.ToLower(v)]
			if ok {
				equivalentFieldId := customFieldMap[strings.ToLower(v)]
				substitutedfieldsToRetrieve = append(substitutedfieldsToRetrieve, fmt.Sprint(equivalentFieldId))
			} else {
				substitutedfieldsToRetrieve = append(substitutedfieldsToRetrieve, strings.ToLower(v))
			}
		}
	}
	processedValues := ProcessedData{filters: substitutedFilters, fields: substitutedfieldsToRetrieve}
	customFieldProcessorChannel <- processedValues
}

// InitiateSearch initiates the search process
func InitiateSearch(config Configuration, customFieldProcessorChannel chan ProcessedData, issueRetrievedChannel chan JiraIssue, totalRestCalls *int, communicator Communicator) {
	processedValues, ok := <-customFieldProcessorChannel
	if ok {
		*totalRestCalls++
		go SearchIssues(config, GetJql(processedValues.filters), processedValues.fields, issueRetrievedChannel, communicator)
	}
}
