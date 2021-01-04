package search

import (
	"encoding/csv"
	"fmt"

	"github.com/pkg/errors"

	"os"
	"strings"
)

func export(results [][]string) error {
	if len(results) == 0 {
		fmt.Printf("No issues found to download")
		return nil
	}

	file, err := os.Create("output.csv")
	if err != nil {
		return errors.Wrapf(err, "failed to create file")
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return errors.Wrapf(writer.WriteAll(results), "failed to write into to csv file")
}

func getJql(filters []Field) string {
	index := 0
	totalCount := len(filters)
	var b strings.Builder
	for _, v := range filters {
		index++
		if strings.Contains(v.Value, ",") {
			valSlice := strings.Split(v.Value, ",")
			b.WriteString(v.ID + " in (" + getInFilterValue(valSlice) + ")")
		} else {
			b.WriteString(v.ID + "=" + "'" + v.Value + "'")
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

// getFieldVal gets the value from the 'fields' property of the issue
func getFieldVal(issue map[string]interface{}, field string) string {
	fieldsMap := issue["fields"].(map[string]interface{})
	val := getStrVal(fieldsMap[field], field)

	//Since we are downloading to csv format, any commas in the field values will lead to improper download
	//So replace all the commas with spaces
	return strings.Replace(val, ",", "", -1)
}

// getStrVal gets the string value based on the type of interface
func getStrVal(field interface{}, fieldName string) string {
	result := "N/A"

	//See if the field is a plain string field or a struct by itself, if it is a struct we need to display one of it's attributes eg: name,displayName etc
	val, ok := field.(map[string]interface{})
	if ok {
		tmp, ok := val[getAttr(fieldName)]
		if ok {
			result = tmp.(string)
		}
	} else if val != nil {
		result = fmt.Sprint(val)
	}

	return result
}

// getAttr gets the attribute of the field, eg name, displayName etc.
func getAttr(fieldName string) string {

	if strings.ToLower(fieldName) == "assignee" || strings.ToLower(fieldName) == "reporter" {
		return "displayName"
	}

	if strings.ToLower(fieldName) == "sprint" || strings.ToLower(fieldName) == "issuetype" || strings.ToLower(fieldName) == "status" || strings.ToLower(fieldName) == "priority" {
		return "name"
	}

	if strings.ToLower(fieldName) == "timetracking" {
		return "originalEstimate"
	}

	if strings.Contains(fieldName, "customfield_") {
		return "name"
	}

	return "value"
}

func getFieldNames(fields []Field) []string {
	names := make([]string, 0)
	for _, f := range fields {
		names = append(names, f.Name)
	}

	return names
}

func getFieldIDs(fields []Field) []string {
	IDs := make([]string, 0)
	for _, f := range fields {
		IDs = append(IDs, f.ID)
	}

	return IDs

}

//HandleError handles errors
func HandleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
