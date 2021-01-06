package search

//Writer writes the result to wherever required.
type Writer interface {
	Write(r Result, p Param) error
}

//CsvWriter writes the results to a csv file
type CsvWriter struct {
	outFile string
}

//NewCsvWriter gives a new csv writer which writes to the outFile file name
func NewCsvWriter(outFile string) CsvWriter {
	return CsvWriter{outFile: outFile}
}

//Write the result from jira to the configured file
func (w CsvWriter) Write(r Result, p Param) error {

	fieldNames, fieldIDs := getFieldNames(p.Fields), getFieldIDs(p.Fields)
	//Write field names to the header of csv
	output := [][]string{fieldNames}

	for _, issue := range r.Data {
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
