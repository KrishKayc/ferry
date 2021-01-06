package search

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	httprequest "github.com/gojira/jira/httprequest"
)

//Result represents the issue results from JIRA rest api
type Result struct {
	StartAt    int                      `json:"startAt"`
	MaxResults int                      `json:"maxResults"`
	Total      int                      `json:"total"`
	Data       []map[string]interface{} `json:"issues"`
}

var client httprequest.Client

const (
	start int = 0
	step  int = 1
)

func debug(data interface{}) {
	file, _ := json.Marshal(data)
	_ = ioutil.WriteFile("debug.json", file, 0644)
}

//Search finds the issue from jira based on the config
func Search(p Param, c httprequest.Client) (Result, error) {
	client = c
	//get all the fields available in the server to set the IDs for custom fields to create the jql query
	r, err := fields()
	if err != nil {
		return r, err
	}
	for _, field := range r.Data {
		setID(p.Filters, field, true)
		setID(p.Fields, field, false)
	}

	return search(getParams(p), make([]map[string]interface{}, 0), start, step)
}

func setID(fields []Field, field map[string]interface{}, isFilter bool) {
	for i, v := range fields {
		if strings.ToLower(field["name"].(string)) == strings.ToLower(v.Name) {
			id := strings.ToLower(strings.ReplaceAll(v.Name, " ", ""))
			if field["custom"].(bool) {
				//For filters, in jql we need to give in the format cf[100203] and for fields to retrieve we need to provide the exact id
				if isFilter {
					id = "cf[" + strings.Replace(field["id"].(string), "customfield_", "", -1) + "]"
				} else {
					id = fmt.Sprint(field["id"].(string))
				}
			}
			v.ID = id
			fields[i] = v
		}
	}
}

func getParams(p Param) map[string]string {
	params := make(map[string]string)
	params["jql"] = getJql(p.Filters)
	params["fields"] = strings.Join(getFieldIDs(p.Fields), ",")
	return params

}

func search(params map[string]string, data []map[string]interface{}, start int, step int) (Result, error) {
	params["startAt"], params["maxResults"] = strconv.Itoa(start), strconv.Itoa(step)
	result, err := searchP(params)
	if err != nil {
		return result, err
	}
	result.Data = append(result.Data, data...)
	if result.Total <= len(result.Data) {
		return result, nil
	}

	return search(params, result.Data, start+step, step)

}

func fields() (Result, error) {
	body := client.Get("/rest/api/2/field", nil)

	var fields []map[string]interface{}
	if err := json.Unmarshal(body, &fields); err != nil {
		return Result{}, err
	}
	return Result{Data: fields}, nil

}

func searchP(params map[string]string) (Result, error) {
	result := Result{}

	body := client.Get("/rest/api/2/search", params)

	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}
	return result, nil
}
