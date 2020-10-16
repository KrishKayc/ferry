package jirafinder

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"

	"os"
	"strconv"
	"strings"
	"time"
)

func writeToCsv(results [][]string, path string) error {
	if len(results) == 0 {
		fmt.Printf("No issues found to download")
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "failed to create file")
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return errors.Wrapf(writer.WriteAll(results), "failed to write into to csv file")
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
		if subTask.TaskType == "Functional Bug" {
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
	if issue == nil {
		return ""
	}
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

//HandleError handles errors
func HandleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func clean(filters map[string]string) {
	for k1, v1 := range filters {
		for k2, v2 := range filters {
			if v1 == v2 && k1 != k2 {
				if strings.HasPrefix(k1, "cf[") {
					delete(filters, k1)
				}
				if strings.HasPrefix(k2, "cf[") {
					delete(filters, k2)
				}
			}
		}
	}
}
