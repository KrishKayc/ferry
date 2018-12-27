package main

import(
	"fmt"
	"encoding/csv"
	"os"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
	 ui "github.com/gizak/termui"
)

// HandleError panics when there is an error in the App
func HandleError(err error){
     if(err != nil){
		panic(err.Error())
	 }
}


// Encodes the value to a Base64 string
func EncodeStringToBase64(val string)string{
	return base64.StdEncoding.EncodeToString([]byte(val))
}


// Reads value from configuration file and return struct
func ReadConfigFromFile(fileName string)Configuration{
	var config Configuration
	fmt.Println("Fetching data based on the configuration file => "+ fileName)
	jsonFile, err := os.Open(fileName)
	HandleError(err)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &config)

	config.AuthToken = EncodeStringToBase64(config.Credentials.Username+":"+ config.Credentials.Password)

	return config
}


// Writes the results to the path specified as csv
func WriteToCsv(results [][]string, path string){

	if(len(results) > 0){

		file, err := os.Create(path)

	    HandleError(err)

		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()
		
		err = writer.WriteAll(results)

		HandleError(err)

	}else{
		fmt.Println("No issues found to download")
	}
}


// Gets the field value based on the field name
func GetFieldValue(field string, issue JiraIssue)string{
	if(field == "assignee"){
		if(issue.AssigneeName != ""){
			return issue.AssigneeName
		}
	   return GetDevTaskAssigneeName(issue.SubTasks)
	 }else if(field == "bug count"){
	   return fmt.Sprint(GetNumberOfFunctionalBugs(issue.SubTasks))
	 }else if(field == "complexity"){
	   return GetComplexityBasedOnDevEstimation(issue.SubTasks)
	 }else{
	   return GetValueFromField(issue.Data, field)
	 } 
}


// Gets Assignee name of the dev task, exclude code review task
func GetDevTaskAssigneeName(subTasks []SubTask)string{
	for _,subTask := range subTasks{
		if(strings.Contains(subTask.Name,"Dev") && !strings.Contains(subTask.Name, "code review")){
			return subTask.AssigneeName
		}
	}

	return "N/A"
}


// Gets the total number of functional issues in the sub tasks
func GetNumberOfFunctionalBugs(subTasks []SubTask)int{
  numberOfFunctionalBugs := 0
  for _,subTask := range subTasks{
	  if(subTask.Type == "Functional Bug"){
		  numberOfFunctionalBugs++
	  }
  }
  return numberOfFunctionalBugs
}


// Gets the complexity based on dev estimation
func GetComplexityBasedOnDevEstimation(subTasks []SubTask)string{
  totalHours := 0
  for _,subTask := range subTasks{
	  if(strings.Contains(subTask.Name,"Dev") && !strings.Contains(subTask.Name, "code review")){
		  hours, _ := strconv.Atoi(strings.TrimRight(subTask.TotalHours,"h"))
		  totalHours += hours
	  }
  }

  if(totalHours <= 8){
	  return "Extra Small"
  }else if(totalHours >= 9 && totalHours <= 16){
	  return "Small"
  }else if(totalHours >= 17 && totalHours <= 24){
	  return "Medium"
  }else if(totalHours >= 25 && totalHours <= 32){
	  return "Large"
  }else if(totalHours >= 33){
	  return "Complex"
  }
  return "N/A"
}


// Displays the Download Progress, total issue count, api calls and time taken in the output terminal
func DisplayProgressAndStatistics(totalIssueCount int, currentIssueCount int,totalApiCalls int, totalTime int, g *ui.Gauge,bc *ui.BarChart){

	bc.Data = []int{totalIssueCount, totalApiCalls, totalTime}

	var percentage int
	progress := ((totalIssueCount - currentIssueCount) % totalIssueCount)

	if(currentIssueCount == -1){
		percentage = 0
	}else if(currentIssueCount == 0){
		percentage = 100
	}else{
		percentage = int(100.0 / (float64(totalIssueCount) / float64(progress)))
	}
	
	g.Percent = percentage
	ui.Render(g, bc)
}


// Gets the download progress and output bar for displaying statistics in terminal
func GetProgressAndStatisticsBar()(*ui.Gauge,*ui.BarChart){
	
	g := ui.NewGauge()
	g.Percent = 0
	g.Width = 50
	g.Height = 3
	g.X = 5
	g.Y = 3
	g.BorderLabel = "Download Progress"
	g.BarColor = ui.ColorRed
	g.BorderFg = ui.ColorWhite
	g.BorderLabelFg = ui.ColorCyan

	barchartData := []int{0, 0, 0}
	bc := ui.NewBarChart()
	bc.BorderLabel = "Statistics"
	bc.Data = barchartData
	bc.Width = 50
	bc.Height = 22
	bc.X = 5
	bc.Y = 7
	bc.BarGap = 1
	bc.BarWidth = 14
	bc.DataLabels = []string{"Issues", "Calls To Jira", "TimeTaken (s)"}
	bc.BarColor = ui.ColorGreen
	bc.NumColor = ui.ColorBlack

	ui.Render(g, bc)

	return g, bc
}


