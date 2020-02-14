package jirafinder

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"

	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func writeToCsv(results [][]string, path string) {

	if len(results) > 0 {

		file, err := os.Create(path)

		HandleError(err)

		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		err = writer.WriteAll(results)

		HandleError(err)

	} else {
		fmt.Println("No issues found to download")
	}
}

func getJql(filters map[string]string) string {
	index := 0
	totalCount := len(filters)
	var b strings.Builder
	for k, v := range filters {
		index++
		if strings.Contains(v, ",") {
			valSlice := strings.Split(v, ",")
			b.WriteString(k + " in (" + getInFilterValue(valSlice) + ")")
		} else {
			b.WriteString(k + "=" + "'" + v + "'")
		}

		if index != totalCount {
			b.WriteString(" AND ")
		}

	}

	return b.String()
}

func getInFilterValue(values []string) string {
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

// GetFieldValue gets the field value based on the field name
func getFieldValue(field string, issue JiraIssue) string {
	if field == "assignee" {
		if issue.AssigneeName != "" {
			return issue.AssigneeName
		}
		return getDevTaskAssigneeName(issue.SubTasks)
	} else if field == "bug count" {
		return fmt.Sprint(getNumberOfFunctionalBugs(issue.SubTasks))
	} else if field == "complexity" {
		return getComplexityBasedOnDevEstimation(issue.SubTasks)
	}

	return getValueFromField(issue.Data, field)
}

// GetValueFromField gets the value from the 'fields' property of the issue
func getValueFromField(issue map[string]interface{}, field string) string {
	val, ok := issue["fields"]
	if ok {
		fieldsMap := val.(map[string]interface{})

		val, ok := fieldsMap[field]
		if ok {
			if strings.ToLower(field) == "created" {
				dateVal, _ := time.Parse("2006-01-02T15:04:05.999-0700", val.(string))
				return dateVal.Format("02/Jan/06")
			}
			return strings.Replace(getValue(val, field), ",", "", -1)
		}
	}
	return "N/A"
}

// GetValue gets the value based on the type of interface
func getValue(val interface{}, fieldName string) string {
	var result string
	arrayVal, isArray := val.([]interface{})
	mapVal, isMap := val.(map[string]interface{})
	if isArray {
		result = arrayVal[0].(map[string]interface{})["value"].(string)
	} else if isMap {
		tmpResult, ok := mapVal[getNestedMapKeyName(fieldName)]
		if ok {
			result = tmpResult.(string)
		}
	} else if val != nil {
		result = fmt.Sprint(val)
	}

	return result
}

// GetNestedMapKeyName gets the nested field name to search for a parent name
func getNestedMapKeyName(fieldName string) string {
	if strings.ToLower(fieldName) == "assignee" || strings.ToLower(fieldName) == "reporter" {
		return "displayName"
	}

	if strings.ToLower(fieldName) == "issuetype" || strings.ToLower(fieldName) == "status" || strings.ToLower(fieldName) == "priority" {
		return "name"
	}

	if strings.ToLower(fieldName) == "timetracking" {
		return "originalEstimate"
	}

	return "value"
}

// GetDevTaskAssigneeName gets Assignee name of the dev task, exclude code review task
func getDevTaskAssigneeName(subTasks []SubTask) string {
	for _, subTask := range subTasks {
		if strings.Contains(subTask.Name, "Dev") && !strings.Contains(subTask.Name, "code review") {
			return subTask.AssigneeName
		}
	}

	return "N/A"
}

// GetNumberOfFunctionalBugs gets the total number of functional issues in the sub tasks
func getNumberOfFunctionalBugs(subTasks []SubTask) int {
	numberOfFunctionalBugs := 0
	for _, subTask := range subTasks {
		if subTask.Type == "Functional Bug" {
			numberOfFunctionalBugs++
		}
	}
	return numberOfFunctionalBugs
}

// GetComplexityBasedOnDevEstimation gets the complexity based on dev estimation
func getComplexityBasedOnDevEstimation(subTasks []SubTask) string {
	totalHours := 0
	for _, subTask := range subTasks {
		if strings.Contains(subTask.Name, "Dev") && !strings.Contains(subTask.Name, "code review") {
			hours, _ := strconv.Atoi(strings.TrimRight(subTask.TotalHours, "h"))
			totalHours += hours
		}
	}

	if totalHours <= 8 {
		return "Extra Small"
	} else if totalHours >= 9 && totalHours <= 16 {
		return "Small"
	} else if totalHours >= 17 && totalHours <= 24 {
		return "Medium"
	} else if totalHours >= 25 && totalHours <= 32 {
		return "Large"
	} else if totalHours >= 33 {
		return "Complex"
	}
	return "N/A"
}

func isBug(issueType string) bool {
	return strings.ToLower(issueType) == "bug" || strings.ToLower(issueType) == "functional bug" || strings.ToLower(issueType) == "production issue"
}

func getDeveloperNameFromLog(issue map[string]interface{}) string {
	developerName := ""
	histories := issue["changelog"].(map[string]interface{})["histories"].([]interface{})
	for _, history := range histories {
		mapHistory := history.(map[string]interface{})
		items := mapHistory["items"].([]interface{})
		for _, item := range items {
			strInDevelopment, ok := item.(map[string]interface{})["toString"].(string)
			if ok && strInDevelopment == "In Development" {
				developerName = mapHistory["author"].(map[string]interface{})["displayName"].(string)
				break
			}
		}

		if developerName != "" {
			break
		}
	}

	return developerName

}

func encodeStringToBase64(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}

// DisplayProgressAndStatistics displays the Download Progress, total issue count, api calls and time taken in the output terminal
func displayProgressAndStatistics(totalIssueCount int, currentIssueCount int, totalAPICalls int, totalTime int) {

	fmt.Print("\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r\r")
	red := color.New(color.FgHiRed).PrintFunc()
	cyan := color.New(color.FgHiCyan).PrintFunc()
	blue := color.New(color.FgHiBlue).PrintFunc()
	yellow := color.New(color.FgHiYellow).PrintFunc()

	cyan(" Total issues : " + strconv.Itoa(totalIssueCount))
	fmt.Print(" | ")
	red(" Total API Calls : " + strconv.Itoa(totalAPICalls))
	fmt.Print(" | ")
	blue(" Total Time : " + strconv.Itoa(totalTime) + " (s)")
	fmt.Print(" | ")

	var percentage int
	progress := ((totalIssueCount - currentIssueCount) % totalIssueCount)

	if currentIssueCount == -1 {
		percentage = 0
	} else if currentIssueCount == 0 {
		percentage = 100
	} else {
		percentage = int(100.0 / (float64(totalIssueCount) / float64(progress)))
	}

	yellow(" Progress : " + strconv.Itoa(percentage) + "%")

}

//HandleError handles errors
func HandleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
