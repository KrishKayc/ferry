package jirafinder

import (
	"encoding/json"
	"fmt"
	"github.com/gojira/ferry/config"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"strings"
	"sync"

	httprequest "github.com/gojira/ferry/httprequest"
)

type keyPairValue struct {
	key   string
	value string
}

type fieldParam struct {
	key  int
	name string
}

type SearchResult struct {
	StartAt    int           `json:"startAt"`
	MaxResults int           `json:"maxResults"`
	Total      int           `json:"total"`
	Issues     []interface{} `json:"issues"`
}

type SubTask struct {
	TaskType     string
	AssigneeName string
	TotalHours   string
	Name         string
}

type JiraIssue struct {
	Data         map[string]interface{}
	SubTasks     []SubTask
	Fields       []string
	AssigneeName string
}

// JiraFinder finds the issue from jira based on the config
type JiraFinder struct {
	Config    config.Configuration
	api       *httprequest.JiraClient
	filtersCh chan keyPairValue
	fieldsCh  chan fieldParam
	fieldKeys []string
	mu        sync.RWMutex
}

func NewJiraFinderFomFile(configFile string) (error, *JiraFinder) {
	err, c := config.New(configFile)
	if err != nil {
		return err, nil
	}

	return NewJiraFinder(c)
}

//NewJiraFinder gives a new jira finder with configurations from the config file
func NewJiraFinder(c *config.Configuration) (error, *JiraFinder) {
	if c.JiraURL == "" {
		return errors.New("no config file found. Set the config first before searching using SetConfig() func"), nil
	}

	return nil, &JiraFinder{
		Config: *c,
		api:    httprequest.NewClient(c.JiraURL, c.AuthToken),

		filtersCh: make(chan keyPairValue),
		fieldsCh:  make(chan fieldParam),

		fieldKeys: make([]string, len(c.FieldsToRetrieve)),
		mu:        sync.RWMutex{},
	}
}

// UseStub enforces usage of httptest
func (f *JiraFinder) UseStub() {
	f.api.UseStub()
}

//Search finds the issue from jira based on the config
func (f *JiraFinder) Search() error {
	output := [][]string{f.Config.FieldsToRetrieve}

	err, out := f.produceFields()
	if err != nil {
		return err
	}

	filters, fields := f.processFields(out)
	err, response := f.search(filters, fields)
	if err != nil {
		return err
	}

	issues := f.prepareIssueObjects(response, fields)
	issueCh := f.processIssues(issues)

	count := 0
	for i := range issueCh {
		if i != nil {
			if f := download(*i); f != nil {
				output = append(output, f)
			}
		}

		count++
		if count == response.Total {
			close(issueCh)
		}
	}

	return writeToCsv(output, f.Config.DownloadPath)
}

func (f *JiraFinder) produceFields() (error, []map[string]interface{}) {
	body := f.api.Get("/rest/api/2/field", nil)

	var fields []map[string]interface{}
	err := json.Unmarshal(body, &fields)
	if err != nil {
		return errors.Wrap(err, "failed to build fields"), nil
	}

	return nil, fields
}

func (f *JiraFinder) collectParams(kpDestination map[string]string) {
	for {
		select {
		case kv, open := <-f.filtersCh:
			if open {
				kpDestination[kv.key] = kv.value
			}

		case fp, open := <-f.fieldsCh:
			if open && fp.name != "" {
				f.addField(fp)
			}
		}
	}
}

func (f *JiraFinder) processFields(fields []map[string]interface{}) (map[string]string, []string) {

	filters := make(map[string]string)

	var wg sync.WaitGroup
	wg.Add(len(fields))

	go f.collectParams(filters)

	for _, field := range fields {
		go func(field map[string]interface{}) {
			defer wg.Done()

			for k, v := range f.Config.Filters {
				if strings.ToLower(field["name"].(string)) == strings.ToLower(k) {
					key := k
					if field["custom"].(bool) {
						key = "cf[" + strings.Replace(field["id"].(string), "customfield_", "", -1) + "]"
					}
					f.filtersCh <- keyPairValue{key, v.(string)}
				}
			}

			for i, v := range f.Config.FieldsToRetrieve {
				if strings.ToLower(field["name"].(string)) == strings.ToLower(v) {
					val := v
					if field["custom"].(bool) {
						val = fmt.Sprint(field["id"].(string))
					}
					f.fieldsCh <- fieldParam{i, val}
				}
			}
		}(field)
	}

	wg.Wait()

	close(f.filtersCh)
	close(f.fieldsCh)
	clean(filters)

	return filters, f.fieldKeys
}

func (f *JiraFinder) addField(field fieldParam) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.fieldKeys[field.key] = field.name
}

func (f *JiraFinder) setFields(params map[string]string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// prevent data race
	params["fields"] = strings.Join(f.fieldKeys, ",")
}

func (f *JiraFinder) search(filters map[string]string, fields []string) (error, *SearchResult) {
	var step int64 = 100
	var startAt int64 = 0
	params := make(map[string]string)
	params["jql"] = getJql(filters)
	params["maxResults"] = strconv.FormatInt(step, 10)
	params["startAt"] = strconv.FormatInt(startAt, 10)
	f.setFields(params)

	err, result := f.doSearchByParams(params)
	if err != nil {
		return err, nil
	}

	// handle results over the limit of 100
	for {
		if result.Total <= len(result.Issues) {
			break
		}

		startAt += step
		params["startAt"] = strconv.FormatInt(startAt, 10)

		err, r := f.doSearchByParams(params)
		if err != nil {
			return err, nil
		}

		result.Issues = append(result.Issues, r.Issues...)
	}

	return nil, result
}

func (f *JiraFinder) doSearchByParams(params map[string]string) (error, *SearchResult) {
	result := new(SearchResult)

	body := f.api.Get("/rest/api/2/search", params)

	if err := json.Unmarshal(body, &result); err != nil {
		return errors.Wrapf(err, "failed to parse search API response"), nil
	}

	return nil, result
}

func (f *JiraFinder) prepareIssueObjects(result *SearchResult, fields []string) []JiraIssue {
	ji := make([]JiraIssue, 0)
	for _, rawIssue := range result.Issues {
		if issue, ok := rawIssue.(map[string]interface{}); ok {
			ji = append(ji, JiraIssue{Data: issue, Fields: fields})
		}
	}

	return ji
}

func (f *JiraFinder) processIssues(issues []JiraIssue) chan *JiraIssue {

	out := make(chan *JiraIssue, 100)
	for i, issue := range issues {
		go func(issue JiraIssue, i int) {
			issueID := issue.Data["id"].(string)
			err, parent := f.getIssue(issueID, true)

			if err != nil {
				log.Printf("error while processing issue %s: %s", issueID, err)
				out <- nil
				return
			}

			subTasks := parent["fields"].(map[string]interface{})["subtasks"].([]interface{})
			result := make([]SubTask, 0)

			for _, v := range subTasks {
				_, subTaskIssue := f.getIssue(v.(map[string]interface{})["id"].(string), false)
				assignee := getValueFromField(subTaskIssue, "assignee")
				issueType := getValueFromField(subTaskIssue, "issuetype")
				name := getValueFromField(subTaskIssue, "summary")
				totalHours := getValueFromField(subTaskIssue, "timetracking")
				currentSubTask := SubTask{TaskType: issueType, Name: name, AssigneeName: assignee, TotalHours: totalHours}

				result = append(result, currentSubTask)
			}

			issue.SubTasks = result

			parentIssueType := getValueFromField(parent, "issuetype")
			if isBug(parentIssueType) {
				issue.AssigneeName = getDeveloperNameFromLog(parent)
			}
			out <- &issue
		}(issue, i)
	}

	return out
}

func (f *JiraFinder) getIssue(issueID string, includeChangeLog bool) (error, map[string]interface{}) {
	var responseResult map[string]interface{}
	var getIssueURL string

	getIssueURL = "/rest/api/2/issue/" + issueID

	if includeChangeLog {
		getIssueURL += "?expand=changelog"
	}

	body := f.api.Get(getIssueURL, nil)

	if err := json.Unmarshal(body, &responseResult); err != nil {
		return errors.Wrapf(err, "failed to retrieve issue"), responseResult
	}

	return nil, responseResult
}

func download(issue JiraIssue) []string {
	fieldValues := make([]string, 0)

	// Listen to final populated issue and prepare the output for all the fields mentioned in the configuration
	for _, field := range issue.Fields {
		val, ok := issue.Data[field]
		if ok {
			fieldValues = append(fieldValues, strings.Replace(val.(string), ",", "", -1))
		} else {
			fieldValues = append(fieldValues, getFieldValue(field, issue))
		}
	}
	if len(fieldValues) > 0 {
		return fieldValues
	}

	return []string{}
}
