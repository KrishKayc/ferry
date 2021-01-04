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
	Err        error
}

func debug(data interface{}) {
	file, _ := json.Marshal(data)
	_ = ioutil.WriteFile("debug.json", file, 0644)
}

//Search finds the issue from jira based on the config
func Search(p Param, client httprequest.Client) error {

	//get all the fields available in the server
	//this is to get the IDs for custom fields to create the jql query
	r := <-allFields(client)

	if r.Err != nil {
		return r.Err
	}

	setID(p, r.Data)
	response, err := search(client, p)
	if err != nil {
		return err
	}

	return download(response.Data, p)
}

func allFields(client httprequest.Client) chan Result {

	c := make(chan Result)
	go func() {
		body := client.Get("/rest/api/2/field", nil)

		var fields []map[string]interface{}
		err := json.Unmarshal(body, &fields)
		if err != nil {
			c <- Result{Err: err}
			return
		}
		c <- Result{Data: fields}
	}()
	return c
}

func setID(p Param, fields []map[string]interface{}) {

	for _, field := range fields {
		setFilterID(field, p)
		setFieldID(field, p)
	}
}

func setFilterID(field map[string]interface{}, p Param) {
	for i, v := range p.Filters {
		if strings.ToLower(field["name"].(string)) == strings.ToLower(v.Name) {
			id := strings.ToLower(strings.ReplaceAll(v.Name, " ", ""))
			if field["custom"].(bool) {
				id = "cf[" + strings.Replace(field["id"].(string), "customfield_", "", -1) + "]"
			}
			v.ID = id
			p.Filters[i] = v
		}
	}
}

func setFieldID(field map[string]interface{}, p Param) {
	for i, v := range p.Fields {
		if strings.ToLower(field["name"].(string)) == strings.ToLower(v.Name) {
			id := strings.ToLower(strings.ReplaceAll(v.Name, " ", ""))

			if field["custom"].(bool) {
				id = fmt.Sprint(field["id"].(string))
			}
			v.ID = id
			p.Fields[i] = v
		}
	}
}

func search(client httprequest.Client, p Param) (*Result, error) {
	var step int64 = 100
	var startAt int64 = 0
	params := make(map[string]string)
	params["jql"] = getJql(p.Filters)
	params["maxResults"] = strconv.FormatInt(step, 10)
	params["startAt"] = strconv.FormatInt(startAt, 10)
	params["fields"] = strings.Join(getFieldIDs(p.Fields), ",")

	result := <-searchP(client, params)
	if result.Err != nil {
		return nil, result.Err
	}

	// handle results over the limit of 100
	for {
		if result.Total <= len(result.Data) {
			break
		}

		startAt += step
		params["startAt"] = strconv.FormatInt(startAt, 10)

		r := <-searchP(client, params)
		if r.Err != nil {
			return nil, r.Err
		}

		result.Data = append(result.Data, r.Data...)
	}

	return result, nil
}

func searchP(client httprequest.Client, params map[string]string) chan *Result {
	c := make(chan *Result)

	go func() {
		result := new(Result)

		body := client.Get("/rest/api/2/search", params)

		if err := json.Unmarshal(body, &result); err != nil {
			c <- &Result{Err: err}
			return
		}
		c <- result

	}()

	return c
}
