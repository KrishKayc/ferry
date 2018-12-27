package main

import "testing"

func TestGetNumberOfFunctionalBugs(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{Type: "Functional Bug"}
	subTask2 := SubTask{Type: "Story"}

	subTasks = append(subTasks, subTask1, subTask2)

	numberOfFunctionalBugs := GetNumberOfFunctionalBugs(subTasks)

	if numberOfFunctionalBugs != 1 {
		t.Errorf("The number of funcational bugs is incorrect, got : %d, want : %d", numberOfFunctionalBugs, 1)
	}
}

func TestGetDevTaskAssigneeNameReturnsDevName(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{Name: "Dev : Analysis", AssigneeName: "Dev1"}

	subTasks = append(subTasks, subTask1)

	devTaskAssigneeName := GetDevTaskAssigneeName(subTasks)

	if devTaskAssigneeName != "Dev1" {
		t.Errorf("The dev task assignee name is wrong, got : %s, want : %s", devTaskAssigneeName, "Dev1")
	}
}

func TestGetDevTaskAssigneeNameNotReturnsReviewerName(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{Name: "Dev : Do code review", AssigneeName: "Dev1"}
	subTask2 := SubTask{Name: "Dev : Coding", AssigneeName: "Dev2"}

	subTasks = append(subTasks, subTask1, subTask2)

	devTaskAssigneeName := GetDevTaskAssigneeName(subTasks)

	if devTaskAssigneeName != "Dev2" {
		t.Errorf("The dev task assignee name is wrong, got : %s, want : %s", devTaskAssigneeName, "Dev2")
	}
}

func TestGetDevTaskAssigneeNameEmpty(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{Name: "QA : Testing", AssigneeName: "Dev1"}
	subTask2 := SubTask{Name: "UX : Review", AssigneeName: "Dev2"}

	subTasks = append(subTasks, subTask1, subTask2)

	devTaskAssigneeName := GetDevTaskAssigneeName(subTasks)

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

	complexity := GetComplexityBasedOnDevEstimation(subTasks)

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

	complexity := GetComplexityBasedOnDevEstimation(subTasks)

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

	complexity := GetComplexityBasedOnDevEstimation(subTasks)

	if complexity != "Large" {
		t.Errorf("Complexity calculation is wrong. got : %s, want : %s", complexity, "Large")
	}
}

func TestGetFieldValueAssigneeFromIssue(t *testing.T) {
	issue := JiraIssue{AssigneeName: "Dev1"}
	fieldValue := GetFieldValue("assignee", issue)

	if fieldValue != "Dev1" {
		t.Errorf("Wrong assignee name, got : %s, want: %s", fieldValue, "Dev1")
	}
}

func TestGetFieldValueAssigneeFromSubTasks(t *testing.T) {
	subTasks := make([]SubTask, 0)
	subTask1 := SubTask{AssigneeName: "Dev1", Name: "Dev : coding"}

	subTasks = append(subTasks, subTask1)
	issue := JiraIssue{SubTasks: subTasks}
	fieldValue := GetFieldValue("assignee", issue)

	if fieldValue != "Dev1" {
		t.Errorf("Wrong assignee name, got : %s, want: %s", fieldValue, "Dev1")
	}
}

func TestGetFieldValueFromField(t *testing.T) {
	issueMap := make(map[string]interface{}, 0)
	fieldsMap := make(map[string]interface{}, 0)

	fieldsMap["story points"] = 5
	issueMap["fields"] = fieldsMap

	fieldValue := GetValueFromField(issueMap, "story points")

	if fieldValue != "5" {
		t.Errorf("Wrong field value. got : %s, want : %s", fieldValue, "5")
	}
}
