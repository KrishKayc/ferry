package search

type Writer interface {
	Write(r Result, p Param) error
}

type CsvWriter struct {
	outFile string
}

func NewCsvWriter(outFile string) CsvWriter {
	return CsvWriter{outFile: outFile}
}

//Download the issues to 'output.csv' file
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
