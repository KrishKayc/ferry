package jirafinder

import (
	"testing"
)

type MockHTTPRequest struct {
}

func (mc *MockHTTPRequest) Send() []byte {
	var jsonString = "{\"fields\" : { \"subTasks\" : [ {} ], \"issuetype\": { \"name\" : \"story\" } }}"
	b := []byte(jsonString)
	return b
}

func TestGetNumberOfFunctionalBugs(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{Type: "Functional Bug"}
	subTask2 := SubTask{Type: "Story"}

	subTasks = append(subTasks, subTask1, subTask2)

	numberOfFunctionalBugs := getNumberOfFunctionalBugs(subTasks)

	if numberOfFunctionalBugs != 1 {
		t.Errorf("The number of funcational bugs is incorrect, got : %d, want : %d", numberOfFunctionalBugs, 1)
	}
}

func TestGetDevTaskAssigneeNameReturnsDevName(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{Name: "Dev : Analysis", AssigneeName: "Dev1"}

	subTasks = append(subTasks, subTask1)

	devTaskAssigneeName := getDevTaskAssigneeName(subTasks)

	if devTaskAssigneeName != "Dev1" {
		t.Errorf("The dev task assignee name is wrong, got : %s, want : %s", devTaskAssigneeName, "Dev1")
	}
}

func TestGetDevTaskAssigneeNameNotReturnsReviewerName(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{Name: "Dev : Do code review", AssigneeName: "Dev1"}
	subTask2 := SubTask{Name: "Dev : Coding", AssigneeName: "Dev2"}

	subTasks = append(subTasks, subTask1, subTask2)

	devTaskAssigneeName := getDevTaskAssigneeName(subTasks)

	if devTaskAssigneeName != "Dev2" {
		t.Errorf("The dev task assignee name is wrong, got : %s, want : %s", devTaskAssigneeName, "Dev2")
	}
}

func TestGetDevTaskAssigneeNameEmpty(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{Name: "QA : Testing", AssigneeName: "Dev1"}
	subTask2 := SubTask{Name: "UX : Review", AssigneeName: "Dev2"}

	subTasks = append(subTasks, subTask1, subTask2)

	devTaskAssigneeName := getDevTaskAssigneeName(subTasks)

	if devTaskAssigneeName != "N/A" {
		t.Errorf("The dev task assignee name is wrong, got : %s, want : %s", devTaskAssigneeName, "N/A")
	}
}

func TestComplexityBasedOnDevEstimates(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{Name: "Dev : Analysis", TotalHours: "8h"}
	subTask2 := SubTask{Name: "Dev : Coding", TotalHours: "12h"}
	subTask3 := SubTask{Name: "Dev : UnitTesting", TotalHours: "12h"}

	subTasks = append(subTasks, subTask1, subTask2, subTask3)

	complexity := getComplexityBasedOnDevEstimation(subTasks)

	if complexity != "Large" {
		t.Errorf("Complexity calculation is wrong. got : %s, want : %s", complexity, "Large")
	}
}

func TestComplexityBasedOnDevEstimatesNotIncludesQATask(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{Name: "Dev : Analysis", TotalHours: "8h"}
	subTask2 := SubTask{Name: "Dev : Coding", TotalHours: "12h"}
	subTask3 := SubTask{Name: "Dev : UnitTesting", TotalHours: "12h"}
	subTask4 := SubTask{Name: "QA : Testing", TotalHours: "12h"}

	subTasks = append(subTasks, subTask1, subTask2, subTask3, subTask4)

	complexity := getComplexityBasedOnDevEstimation(subTasks)

	if complexity != "Large" {
		t.Errorf("Complexity calculation is wrong. got : %s, want : %s", complexity, "Large")
	}
}

func TestComplexityBasedOnDevEstimatesNotIncludesReviewTask(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{Name: "Dev : Analysis", TotalHours: "8h"}
	subTask2 := SubTask{Name: "Dev : Coding", TotalHours: "12h"}
	subTask3 := SubTask{Name: "Dev : UnitTesting", TotalHours: "12h"}
	subTask4 := SubTask{Name: "Dev : code review", TotalHours: "12h"}

	subTasks = append(subTasks, subTask1, subTask2, subTask3, subTask4)

	complexity := getComplexityBasedOnDevEstimation(subTasks)

	if complexity != "Large" {
		t.Errorf("Complexity calculation is wrong. got : %s, want : %s", complexity, "Large")
	}
}

func TestGetFieldValueAssigneeFromIssue(t *testing.T) {
	issue := JiraIssue{AssigneeName: "Dev1"}
	fieldValue := getFieldValue("assignee", issue)

	if fieldValue != "Dev1" {
		t.Errorf("Wrong assignee name, got : %s, want: %s", fieldValue, "Dev1")
	}
}

func TestGetFieldValueAssigneeFromSubTasks(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{AssigneeName: "Dev1", Name: "Dev : coding"}

	subTasks = append(subTasks, subTask1)
	issue := JiraIssue{SubTasks: subTasks}
	fieldValue := getFieldValue("assignee", issue)

	if fieldValue != "Dev1" {
		t.Errorf("Wrong assignee name, got : %s, want: %s", fieldValue, "Dev1")
	}
}

func TestGetFieldValueFromField(t *testing.T) {
	issueMap := make(map[string]interface{}, 0)
	fieldsMap := make(map[string]interface{}, 0)

	fieldsMap["story points"] = 5
	issueMap["fields"] = fieldsMap

	fieldValue := getValueFromField(issueMap, "story points")

	if fieldValue != "5" {
		t.Errorf("Wrong field value. got : %s, want : %s", fieldValue, "5")
	}
}

func TestGetNestedMapKeyName(t *testing.T) {
	result := getNestedMapKeyName("Assignee")

	if result != "displayName" {
		ThrowError(t, "worng assignee nested name", "displayName", result)
	}

	result = getNestedMapKeyName("Reporter")

	if result != "displayName" {
		ThrowError(t, "worng reporter nested name", "displayName", result)
	}

	result = getNestedMapKeyName("IssueType")

	if result != "name" {
		ThrowError(t, "worng IssueType nested name", "name", result)
	}

	result = getNestedMapKeyName("TimeTracking")

	if result != "originalEstimate" {
		ThrowError(t, "worng TimeTracking nested name", "originalEstimate", result)
	}

	result = getNestedMapKeyName("Status")

	if result != "name" {
		ThrowError(t, "worng Status nested name", "name", result)
	}

	result = getNestedMapKeyName("someCustomField")

	if result != "value" {
		ThrowError(t, "worng custom field value name", "value", result)
	}
}

// func TestGetIssue(t *testing.T) {
// 	mc := MockCommunicator{}
// 	issue := getIssue(Configuration{}, "", false, &mc)

// 	if issue["fields"] == nil {
// 		t.Errorf("Get Issue not returning map")
// 	}
// }

// func TestGetIssueIncludingLog(t *testing.T) {
// 	mc := MockCommunicator{}
// 	issue := getIssue(Configuration{}, "", true, &mc)

// 	if issue["fields"] == nil {
// 		t.Errorf("Get Issue not returning map")
// 	}
// }

// func TestGetCustomFields(t *testing.T) {
// 	customFieldChannel := make(chan map[string]string)
// 	mc := MockCommunicatorFieldRetrieval{}

// 	go getCustomFields(Configuration{}, customFieldChannel, &mc)

// 	customFieldMap := <-customFieldChannel

// 	if customFieldMap["testname"] != "testid" {
// 		t.Errorf("Failed retieving custom field values into a map")
// 	}
// }

// func TestGetCustomFieldsStaticFields(t *testing.T) {
// 	customFieldChannel := make(chan map[string]string)
// 	mc := MockCommunicatorFieldRetrievalStaticFields{}

// 	go GetCustomFields(Configuration{}, customFieldChannel, &mc)

// 	customFieldMap := <-customFieldChannel

// 	if len(customFieldMap) != 0 {
// 		t.Errorf("Failed retieving custom field values into a map")
// 	}
// }
func ThrowError(t *testing.T, errorMsg string, expected string, actual string) {
	t.Errorf("%s, got : %s, want: %s", errorMsg, actual, expected)
}
