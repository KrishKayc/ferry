package search

import (
	"encoding/csv"
	"fmt"

	"github.com/pkg/errors"

	"os"
	"strings"
)

//This specifies which attributes must be displayed to the user in the exported csv.
//Eg : if field 'assignee' is provided by the user in the --fields flag, display 'displayName' and so..
//Add (or) modify this map as per the requirement
var fieldDisplayAttr = map[string]string{
	"assignee":     "displayName",
	"reporter":     "displayName",
	"issuetype":    "name",
	"status":       "name",
	"priority":     "name",
	"timetracking": "originalEstimate",
}

func download(issues []map[string]interface{}, p Param) error {

	fieldNames, fieldIDs := getFieldNames(p.Fields), getFieldIDs(p.Fields)
	//Write field names to the header of csv
	output := [][]string{fieldNames}

	for _, issue := range issues {
		fieldValues := make([]string, 0)

		for _, field := range fieldIDs {
			fieldValues = append(fieldValues, getFieldVal(issue, field))
		}
		if len(fieldValues) > 0 {
			output = append(output, fieldValues)
		}
	}
	return export(output)
}
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
	index, total := 0, len(filters)
	var b strings.Builder
	for _, v := range filters {
		index++
		if strings.Contains(v.Value, ",") {
			arr := strings.Split(v.Value, ",")
			b.WriteString(v.ID + " in (" + getInFilterVal(arr) + ")")
		} else {
			b.WriteString(v.ID + "=" + "'" + v.Value + "'")
		}

		if index != total {
			b.WriteString(" AND ")
		}

	}

	return b.String()
}

func getInFilterVal(values []string) string {
	index, total := 0, len(values)
	var b strings.Builder
	for _, val := range values {
		index++
		b.WriteString("'" + strings.TrimSpace(val) + "'")
		if index != total {
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
	//If the issue does not have the required field, return as N/A.
	if field == nil {
		return "N/A"
	}
	//The field type can be anything. It can be a map or array or plain string.
	//If it is a map, we need to display one of it's attributes eg: name,displayName etc
	val, ok := field.(map[string]interface{})
	if ok {
		return getAttrFromMap(val, fieldName)
	}
	//If it is an array, take the first element which is a map, this happens to custom field eg : Sprint
	arr, ok := field.([]interface{})
	if ok {
		return getAttrFromMap(arr[0].(map[string]interface{}), fieldName)
	}

	//This is for plain string fields
	return field.(string)
}
func getAttrFromMap(field map[string]interface{}, fieldName string) string {
	tmp, ok := field[getAttr(fieldName)]
	if ok {
		return tmp.(string)
	}
	return "N/A"
}

// getAttr gets the attribute of the field, eg name, displayName etc.
func getAttr(fieldName string) string {
	if strings.Contains(fieldName, "customfield_") {
		return "name"
	}

	if val, ok := fieldDisplayAttr[fieldName]; ok {
		return val
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
