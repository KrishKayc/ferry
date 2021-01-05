package main

import (
	"github.com/gojira/jira/cmd"
	"github.com/gojira/jira/search"
)

func main() {
	csvWriter := search.NewCsvWriter("output.csv")
	cmd.Execute(csvWriter)
}
