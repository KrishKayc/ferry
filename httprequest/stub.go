package httprequest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
)

// UseStub is aimed to serve a fake Jira API Rest service
// with same kind of result according to URL requested
func (c *JiraClient) UseStub() {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var resp string

		issueReq, _ := regexp.Compile("/rest/api/2/issue/([0-9]+)(\\?(.*))?$")
		searchReq, _ := regexp.Compile("/rest/api/2/search(\\?(.*))?$")

		switch {
		case r.RequestURI == "/rest/api/2/field":
			resp = `[
  {
    "id": "statuscategorychangedate",
    "key": "statuscategorychangedate",
    "name": "Status Category Changed",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "statusCategoryChangedDate"
    ],
    "schema": {
      "type": "datetime",
      "system": "statuscategorychangedate"
    }
  },
  {
    "id": "issuetype",
    "key": "issuetype",
    "name": "Issue Type",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "issuetype",
      "type"
    ],
    "schema": {
      "type": "issuetype",
      "system": "issuetype"
    }
  },
  {
    "id": "parent",
    "key": "parent",
    "name": "Parent",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": [
      "parent"
    ]
  },
  {
    "id": "timespent",
    "key": "timespent",
    "name": "Time Spent",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": [
      "timespent"
    ],
    "schema": {
      "type": "number",
      "system": "timespent"
    }
  },
  {
    "id": "project",
    "key": "project",
    "name": "Project",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "project"
    ],
    "schema": {
      "type": "project",
      "system": "project"
    }
  },
  {
    "id": "fixVersions",
    "key": "fixVersions",
    "name": "Fix versions",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "fixVersion"
    ],
    "schema": {
      "type": "array",
      "items": "version",
      "system": "fixVersions"
    }
  },
  {
    "id": "aggregatetimespent",
    "key": "aggregatetimespent",
    "name": "Σ Time Spent",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": [],
    "schema": {
      "type": "number",
      "system": "aggregatetimespent"
    }
  },
  {
    "id": "statusCategory",
    "key": "statusCategory",
    "name": "Status Category",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "statusCategory"
    ]
  },
  {
    "id": "resolution",
    "key": "resolution",
    "name": "Resolution",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "resolution"
    ],
    "schema": {
      "type": "resolution",
      "system": "resolution"
    }
  },
  {
    "id": "resolutiondate",
    "key": "resolutiondate",
    "name": "Resolved",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "resolutiondate",
      "resolved"
    ],
    "schema": {
      "type": "datetime",
      "system": "resolutiondate"
    }
  },
  {
    "id": "workratio",
    "key": "workratio",
    "name": "Work Ratio",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "workratio"
    ],
    "schema": {
      "type": "number",
      "system": "workratio"
    }
  },
  {
    "id": "issuerestriction",
    "key": "issuerestriction",
    "name": "Restrict to",
    "custom": false,
    "orderable": true,
    "navigable": false,
    "searchable": true,
    "clauseNames": [],
    "schema": {
      "type": "issuerestriction",
      "system": "issuerestriction"
    }
  },
  {
    "id": "lastViewed",
    "key": "lastViewed",
    "name": "Last Viewed",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": [
      "lastViewed"
    ],
    "schema": {
      "type": "datetime",
      "system": "lastViewed"
    }
  },
  {
    "id": "watches",
    "key": "watches",
    "name": "Watchers",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": [
      "watchers"
    ],
    "schema": {
      "type": "watches",
      "system": "watches"
    }
  },
  {
    "id": "thumbnail",
    "key": "thumbnail",
    "name": "Images",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": []
  },
  {
    "id": "created",
    "key": "created",
    "name": "Created",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "created",
      "createdDate"
    ],
    "schema": {
      "type": "datetime",
      "system": "created"
    }
  },
  {
    "id": "customfield_10020",
    "key": "customfield_10020",
    "name": "Sprint",
    "untranslatedName": "Sprint",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10020]",
      "Sprint"
    ],
    "schema": {
      "type": "array",
      "items": "json",
      "custom": "com.pyxis.greenhopper.jira:gh-sprint",
      "customId": 10020
    }
  },
  {
    "id": "customfield_10021",
    "key": "customfield_10021",
    "name": "Flagged",
    "untranslatedName": "Flagged",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10021]",
      "Flagged",
      "Flagged[Checkboxes]"
    ],
    "schema": {
      "type": "array",
      "items": "option",
      "custom": "com.atlassian.jira.plugin.system.customfieldtypes:multicheckboxes",
      "customId": 10021
    }
  },
  {
    "id": "customfield_10022",
    "key": "customfield_10022",
    "name": "Target start",
    "untranslatedName": "Target start",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10022]",
      "Target start"
    ],
    "schema": {
      "type": "date",
      "custom": "com.atlassian.jpo:jpo-custom-field-baseline-start",
      "customId": 10022
    }
  },
  {
    "id": "priority",
    "key": "priority",
    "name": "Priority",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "priority"
    ],
    "schema": {
      "type": "priority",
      "system": "priority"
    }
  },
  {
    "id": "customfield_10023",
    "key": "customfield_10023",
    "name": "Target end",
    "untranslatedName": "Target end",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10023]",
      "Target end"
    ],
    "schema": {
      "type": "date",
      "custom": "com.atlassian.jpo:jpo-custom-field-baseline-end",
      "customId": 10023
    }
  },
  {
    "id": "customfield_10024",
    "key": "customfield_10024",
    "name": "[CHART] Date of First Response",
    "untranslatedName": "[CHART] Date of First Response",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "[CHART] Date of First Response",
      "[CHART] Date of First Response[Date of first response]",
      "cf[10024]"
    ],
    "schema": {
      "type": "datetime",
      "custom": "com.atlassian.jira.ext.charting:firstresponsedate",
      "customId": 10024
    }
  },
  {
    "id": "customfield_10025",
    "key": "customfield_10025",
    "name": "[CHART] Time in Status",
    "untranslatedName": "[CHART] Time in Status",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "[CHART] Time in Status",
      "[CHART] Time in Status[Time in Status]",
      "cf[10025]"
    ],
    "schema": {
      "type": "any",
      "custom": "com.atlassian.jira.ext.charting:timeinstatus",
      "customId": 10025
    }
  },
  {
    "id": "labels",
    "key": "labels",
    "name": "Labels",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "labels"
    ],
    "schema": {
      "type": "array",
      "items": "string",
      "system": "labels"
    }
  },
  {
    "id": "customfield_10026",
    "key": "customfield_10026",
    "name": "Story Points",
    "untranslatedName": "Story Points",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10026]",
      "Story Points",
      "Story Points[Number]"
    ],
    "schema": {
      "type": "number",
      "custom": "com.atlassian.jira.plugin.system.customfieldtypes:float",
      "customId": 10026
    }
  },
  {
    "id": "customfield_10016",
    "key": "customfield_10016",
    "name": "Story point estimate",
    "untranslatedName": "Story point estimate",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10016]",
      "Story point estimate"
    ],
    "schema": {
      "type": "number",
      "custom": "com.pyxis.greenhopper.jira:jsw-story-points",
      "customId": 10016
    }
  },
  {
    "id": "customfield_10017",
    "key": "customfield_10017",
    "name": "Issue color",
    "untranslatedName": "Issue color",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10017]",
      "Issue color"
    ],
    "schema": {
      "type": "string",
      "custom": "com.pyxis.greenhopper.jira:jsw-issue-color",
      "customId": 10017
    }
  },
  {
    "id": "customfield_10018",
    "key": "customfield_10018",
    "name": "Parent Link",
    "untranslatedName": "Parent Link",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10018]",
      "Parent Link"
    ],
    "schema": {
      "type": "any",
      "custom": "com.atlassian.jpo:jpo-custom-field-parent",
      "customId": 10018
    }
  },
  {
    "id": "customfield_10019",
    "key": "customfield_10019",
    "name": "Rank",
    "untranslatedName": "Rank",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10019]",
      "Rank"
    ],
    "schema": {
      "type": "any",
      "custom": "com.pyxis.greenhopper.jira:gh-lexo-rank",
      "customId": 10019
    }
  },
  {
    "id": "timeestimate",
    "key": "timeestimate",
    "name": "Remaining Estimate",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": [
      "remainingEstimate",
      "timeestimate"
    ],
    "schema": {
      "type": "number",
      "system": "timeestimate"
    }
  },
  {
    "id": "aggregatetimeoriginalestimate",
    "key": "aggregatetimeoriginalestimate",
    "name": "Σ Original Estimate",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": [],
    "schema": {
      "type": "number",
      "system": "aggregatetimeoriginalestimate"
    }
  },
  {
    "id": "versions",
    "key": "versions",
    "name": "Affects versions",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "affectedVersion"
    ],
    "schema": {
      "type": "array",
      "items": "version",
      "system": "versions"
    }
  },
  {
    "id": "issuelinks",
    "key": "issuelinks",
    "name": "Linked Issues",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "issueLink"
    ],
    "schema": {
      "type": "array",
      "items": "issuelinks",
      "system": "issuelinks"
    }
  },
  {
    "id": "assignee",
    "key": "assignee",
    "name": "Assignee",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "assignee"
    ],
    "schema": {
      "type": "user",
      "system": "assignee"
    }
  },
  {
    "id": "updated",
    "key": "updated",
    "name": "Updated",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "updated",
      "updatedDate"
    ],
    "schema": {
      "type": "datetime",
      "system": "updated"
    }
  },
  {
    "id": "status",
    "key": "status",
    "name": "Status",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "status"
    ],
    "schema": {
      "type": "status",
      "system": "status"
    }
  },
  {
    "id": "components",
    "key": "components",
    "name": "Components",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "component"
    ],
    "schema": {
      "type": "array",
      "items": "component",
      "system": "components"
    }
  },
  {
    "id": "issuekey",
    "key": "issuekey",
    "name": "Key",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": [
      "id",
      "issue",
      "issuekey",
      "key"
    ]
  },
  {
    "id": "timeoriginalestimate",
    "key": "timeoriginalestimate",
    "name": "Original Estimate",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": [
      "originalEstimate",
      "timeoriginalestimate"
    ],
    "schema": {
      "type": "number",
      "system": "timeoriginalestimate"
    }
  },
  {
    "id": "description",
    "key": "description",
    "name": "Description",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "description"
    ],
    "schema": {
      "type": "string",
      "system": "description"
    }
  },
  {
    "id": "customfield_10010",
    "key": "customfield_10010",
    "name": "Request Type",
    "untranslatedName": "Request Type",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10010]",
      "Request Type"
    ],
    "schema": {
      "type": "sd-customerrequesttype",
      "custom": "com.atlassian.servicedesk:vp-origin",
      "customId": 10010
    }
  },
  {
    "id": "customfield_10011",
    "key": "customfield_10011",
    "name": "Epic Name",
    "untranslatedName": "Epic Name",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10011]",
      "Epic Name"
    ],
    "schema": {
      "type": "string",
      "custom": "com.pyxis.greenhopper.jira:gh-epic-label",
      "customId": 10011
    }
  },
  {
    "id": "customfield_10012",
    "key": "customfield_10012",
    "name": "Epic Status",
    "untranslatedName": "Epic Status",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10012]",
      "Epic Status"
    ],
    "schema": {
      "type": "option",
      "custom": "com.pyxis.greenhopper.jira:gh-epic-status",
      "customId": 10012
    }
  },
  {
    "id": "customfield_10013",
    "key": "customfield_10013",
    "name": "Epic Color",
    "untranslatedName": "Epic Color",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10013]",
      "Epic Color"
    ],
    "schema": {
      "type": "string",
      "custom": "com.pyxis.greenhopper.jira:gh-epic-color",
      "customId": 10013
    }
  },
  {
    "id": "customfield_10014",
    "key": "customfield_10014",
    "name": "Epic Link",
    "untranslatedName": "Epic Link",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10014]",
      "Epic Link"
    ],
    "schema": {
      "type": "any",
      "custom": "com.pyxis.greenhopper.jira:gh-epic-link",
      "customId": 10014
    }
  },
  {
    "id": "timetracking",
    "key": "timetracking",
    "name": "Time tracking",
    "custom": false,
    "orderable": true,
    "navigable": false,
    "searchable": true,
    "clauseNames": [],
    "schema": {
      "type": "timetracking",
      "system": "timetracking"
    }
  },
  {
    "id": "customfield_10015",
    "key": "customfield_10015",
    "name": "Start date",
    "untranslatedName": "Start date",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10015]",
      "Start date",
      "Start date[Date]"
    ],
    "schema": {
      "type": "date",
      "custom": "com.atlassian.jira.plugin.system.customfieldtypes:datepicker",
      "customId": 10015
    }
  },
  {
    "id": "customfield_10005",
    "key": "customfield_10005",
    "name": "Change type",
    "untranslatedName": "Change type",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10005]",
      "Change type",
      "Change type[Dropdown]"
    ],
    "schema": {
      "type": "option",
      "custom": "com.atlassian.jira.plugin.system.customfieldtypes:select",
      "customId": 10005
    }
  },
  {
    "id": "customfield_10006",
    "key": "customfield_10006",
    "name": "Change risk",
    "untranslatedName": "Change risk",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10006]",
      "Change risk",
      "Change risk[Dropdown]"
    ],
    "schema": {
      "type": "option",
      "custom": "com.atlassian.jira.plugin.system.customfieldtypes:select",
      "customId": 10006
    }
  },
  {
    "id": "security",
    "key": "security",
    "name": "Security Level",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "level"
    ],
    "schema": {
      "type": "securitylevel",
      "system": "security"
    }
  },
  {
    "id": "customfield_10007",
    "key": "customfield_10007",
    "name": "Change reason",
    "untranslatedName": "Change reason",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10007]",
      "Change reason",
      "Change reason[Dropdown]"
    ],
    "schema": {
      "type": "option",
      "custom": "com.atlassian.jira.plugin.system.customfieldtypes:select",
      "customId": 10007
    }
  },
  {
    "id": "customfield_10008",
    "key": "customfield_10008",
    "name": "Change start date",
    "untranslatedName": "Change start date",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10008]",
      "Change start date",
      "Change start date[Time stamp]"
    ],
    "schema": {
      "type": "datetime",
      "custom": "com.atlassian.jira.plugin.system.customfieldtypes:datetime",
      "customId": 10008
    }
  },
  {
    "id": "attachment",
    "key": "attachment",
    "name": "Attachment",
    "custom": false,
    "orderable": true,
    "navigable": false,
    "searchable": true,
    "clauseNames": [
      "attachments"
    ],
    "schema": {
      "type": "array",
      "items": "attachment",
      "system": "attachment"
    }
  },
  {
    "id": "customfield_10009",
    "key": "customfield_10009",
    "name": "Change completion date",
    "untranslatedName": "Change completion date",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10009]",
      "Change completion date",
      "Change completion date[Time stamp]"
    ],
    "schema": {
      "type": "datetime",
      "custom": "com.atlassian.jira.plugin.system.customfieldtypes:datetime",
      "customId": 10009
    }
  },
  {
    "id": "summary",
    "key": "summary",
    "name": "Summary",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "summary"
    ],
    "schema": {
      "type": "string",
      "system": "summary"
    }
  },
  {
    "id": "creator",
    "key": "creator",
    "name": "Creator",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "creator"
    ],
    "schema": {
      "type": "user",
      "system": "creator"
    }
  },
  {
    "id": "subtasks",
    "key": "subtasks",
    "name": "Sub-tasks",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": [
      "subtasks"
    ],
    "schema": {
      "type": "array",
      "items": "issuelinks",
      "system": "subtasks"
    }
  },
  {
    "id": "reporter",
    "key": "reporter",
    "name": "Reporter",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "reporter"
    ],
    "schema": {
      "type": "user",
      "system": "reporter"
    }
  },
  {
    "id": "customfield_10000",
    "key": "customfield_10000",
    "name": "Development",
    "untranslatedName": "development",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10000]",
      "development"
    ],
    "schema": {
      "type": "any",
      "custom": "com.atlassian.jira.plugins.jira-development-integration-plugin:devsummarycf",
      "customId": 10000
    }
  },
  {
    "id": "customfield_10001",
    "key": "customfield_10001",
    "name": "Team",
    "untranslatedName": "Team",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10001]",
      "Team",
      "Team[Team]"
    ],
    "schema": {
      "type": "any",
      "custom": "com.atlassian.teams:rm-teams-custom-field-team",
      "customId": 10001
    }
  },
  {
    "id": "customfield_10002",
    "key": "customfield_10002",
    "name": "Organizations",
    "untranslatedName": "Organizations",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "cf[10002]",
      "Organizations"
    ],
    "schema": {
      "type": "array",
      "items": "sd-customerorganization",
      "custom": "com.atlassian.servicedesk:sd-customer-organizations",
      "customId": 10002
    }
  },
  {
    "id": "customfield_10003",
    "key": "customfield_10003",
    "name": "Approvers",
    "untranslatedName": "Approvers",
    "custom": true,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "Approvers",
      "Approvers[User Picker (multiple users)]",
      "cf[10003]"
    ],
    "schema": {
      "type": "array",
      "items": "user",
      "custom": "com.atlassian.jira.plugin.system.customfieldtypes:multiuserpicker",
      "customId": 10003
    }
  },
  {
    "id": "environment",
    "key": "environment",
    "name": "Environment",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "environment"
    ],
    "schema": {
      "type": "string",
      "system": "environment"
    }
  },
  {
    "id": "duedate",
    "key": "duedate",
    "name": "Due date",
    "custom": false,
    "orderable": true,
    "navigable": true,
    "searchable": true,
    "clauseNames": [
      "due",
      "duedate"
    ],
    "schema": {
      "type": "date",
      "system": "duedate"
    }
  },
  {
    "id": "progress",
    "key": "progress",
    "name": "Progress",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": [
      "progress"
    ],
    "schema": {
      "type": "progress",
      "system": "progress"
    }
  },
  {
    "id": "comment",
    "key": "comment",
    "name": "Comment",
    "custom": false,
    "orderable": true,
    "navigable": false,
    "searchable": true,
    "clauseNames": [
      "comment"
    ],
    "schema": {
      "type": "comments-page",
      "system": "comment"
    }
  },
  {
    "id": "votes",
    "key": "votes",
    "name": "Votes",
    "custom": false,
    "orderable": false,
    "navigable": true,
    "searchable": false,
    "clauseNames": [
      "votes"
    ],
    "schema": {
      "type": "votes",
      "system": "votes"
    }
  },
  {
    "id": "worklog",
    "key": "worklog",
    "name": "Log Work",
    "custom": false,
    "orderable": true,
    "navigable": false,
    "searchable": true,
    "clauseNames": [],
    "schema": {
      "type": "array",
      "items": "worklog",
      "system": "worklog"
    }
  }
]`
		case searchReq.MatchString(r.RequestURI):
			resp = `{
  "expand": "schema,names",
  "startAt": 0,
  "maxResults": 100,
  "total": 6,
  "issues": [
    {
      "expand": "operations,versionedRepresentations,editmeta,changelog,renderedFields",
      "id": "10006",
      "key": "POS-7",
      "fields": {
        "summary": "Reporting",
        "assignee": null,
        "customfield_10026": null
      }
    },
    {
      "expand": "operations,versionedRepresentations,editmeta,changelog,renderedFields",
      "id": "10004",
      "self": "https://myspace.atlassian.net/rest/api/2/issue/10004",
      "key": "POS-5",
      "fields": {
        "summary": "Admin Magasin",
        "assignee": {
          "emailAddress": "user@gmail.com",
          "active": true,
          "accountType": "atlassian"
        },
        "customfield_10026": null
      }
    }
  ]
}`

		case issueReq.MatchString(r.RequestURI):
			m := issueReq.FindStringSubmatch(r.RequestURI)
			issueType := "Story"

			if strings.Contains(r.RequestURI, "expand=changelog") {
				issueType = "Bug"
			}

			resp = fmt.Sprintf(`{
  "expand": "renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations",
  "id": "%s",
  "key": "POS-1",
  "changelog": {
    "startAt": 0,
    "maxResults": 4,
    "total": 4,
    "histories": [
      {
        "id": "10056",
        "author": {
          "emailAddress": "user@gmail.com",
          "displayName": "User Name",
          "active": true,
          "accountType": "atlassian"
        },
        "created": "2020-08-19T20:11:37.133+0300",
        "items": [
          {
            "field": "Sprint",
            "fieldtype": "custom",
            "fieldId": "customfield_10020",
            "from": "",
            "fromString": "",
            "to": "1",
            "toString": "POS Sprint 1"
          }
        ]
      }
    ]
  },
  "fields": {
    "statuscategorychangedate": "2020-08-17T08:13:32.569+0300",
    "issuetype": {
      "id": "10001",
      "description": "Functionality or a feature expressed as a user goal.",
      "name": "%s",
      "subtask": false,
      "avatarId": 10315
    },
    "timespent": null,
    "project": {
      "id": "10000",
      "key": "POS",
      "name": "POS",
      "projectTypeKey": "software",
      "simplified": false
    },
    "fixVersions": [],
    "aggregatetimespent": null,
    "resolution": null,
    "resolutiondate": null,
    "workratio": -1,
    "issuerestriction": {
      "issuerestrictions": {},
      "shouldDisplay": false
    },
    "watches": {
      "watchCount": 1,
      "isWatching": true
    },
    "lastViewed": "2020-08-19T20:11:40.821+0300",
    "created": "2020-08-17T08:13:32.383+0300",
    "customfield_10020": [
      {
        "id": 1,
        "name": "POS Sprint 1",
        "state": "active",
        "boardId": 1,
        "goal": "Implement basic features",
        "startDate": "2020-08-19T17:11:53.299Z",
        "endDate": "2020-09-02T17:11:00.000Z"
      }
    ],
    "customfield_10021": null,
    "customfield_10022": null,
    "priority": {
      "name": "Medium",
      "id": "3"
    },
    "customfield_10023": null,
    "customfield_10024": null,
    "customfield_10025": null,
    "customfield_10026": null,
    "labels": [],
    "customfield_10016": null,
    "customfield_10017": null,
    "customfield_10018": {
      "hasEpicLinkFieldDependency": false,
      "showField": false,
      "nonEditableReason": {
        "reason": "PLUGIN_LICENSE_ERROR",
        "message": "The Parent Link is only available to Jira Premium users."
      }
    },
    "customfield_10019": "0|i0001b:",
    "aggregatetimeoriginalestimate": null,
    "timeestimate": null,
    "versions": [],
    "issuelinks": [],
    "assignee": null,
    "updated": "2020-08-19T20:11:37.130+0300",
    "status": {
      "description": "",
      "name": "To Do",
      "id": "10000",
      "statusCategory": {
        "id": 2,
        "key": "new",
        "colorName": "blue-gray",
        "name": "To Do"
      }
    },
    "components": [],
    "timeoriginalestimate": null,
    "description": "[https://react-material-kit.devias.io/app/reports/dashboard|https://react-material-kit.devias.io/app/reports/dashboard]",
    "customfield_10010": null,
    "customfield_10014": "POS-16",
    "customfield_10015": null,
    "timetracking": {},
    "customfield_10005": null,
    "customfield_10006": null,
    "customfield_10007": null,
    "security": null,
    "customfield_10008": null,
    "aggregatetimeestimate": null,
    "attachment": [],
    "customfield_10009": null,
    "summary": "Dashboard components",
    "creator": {
      "emailAddress": "user@gmail.com",
      "avatarUrls": {
        "48x48": "https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/557058:114f44fd-72f6-409b-9327-a5e61c75fe72/0d402fe1-b810-42d8-9450-3813e764984c/48",
        "24x24": "https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/557058:114f44fd-72f6-409b-9327-a5e61c75fe72/0d402fe1-b810-42d8-9450-3813e764984c/24",
        "16x16": "https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/557058:114f44fd-72f6-409b-9327-a5e61c75fe72/0d402fe1-b810-42d8-9450-3813e764984c/16",
        "32x32": "https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/557058:114f44fd-72f6-409b-9327-a5e61c75fe72/0d402fe1-b810-42d8-9450-3813e764984c/32"
      },
      "displayName": "Jira User",
      "active": true,
      "accountType": "atlassian"
    },
    "subtasks": [
      {
        "id": "10017",
        "key": "POS-18",
        "fields": {
          "summary": "test",
          "status": {
            "description": "",
            "name": "To Do",
            "id": "10000",
            "statusCategory": {
              "id": 2,
              "key": "new",
              "colorName": "blue-gray",
              "name": "To Do"
            }
          },
          "priority": {
            "name": "Medium",
            "id": "3"
          },
          "issuetype": {
            "id": "10003",
            "description": "A small piece of work that's part of a larger task.",
            "name": "Sub-task",
            "subtask": true,
            "avatarId": 10316
          }
        }
      }
    ],
    "reporter": {
      "emailAddress": "user@gmail.com",
      "active": true,
      "accountType": "atlassian"
    },
    "customfield_10000": "{}",
    "aggregateprogress": {
      "progress": 0,
      "total": 0
    },
    "customfield_10001": null,
    "customfield_10002": null,
    "customfield_10003": null,
    "customfield_10004": null,
    "environment": null,
    "duedate": null,
    "progress": {
      "progress": 0,
      "total": 0
    },
    "votes": {
      "votes": 0,
      "hasVoted": false
    },
    "comment": {
      "comments": [],
      "maxResults": 0,
      "total": 0,
      "startAt": 0
    },
    "worklog": {
      "startAt": 0,
      "maxResults": 20,
      "total": 0,
      "worklogs": []
    }
  }
}`, m[1], issueType)

		default:
			resp = `{
  "id": "https://docs.atlassian.com/jira/REST/schema/error-collection#",
  "title": "Error Collection",
  "type": "object",
  "properties": {
    "errorMessages": {
      "type": "array",
      "items": {
        "type": "string"
      }
    },
    "errors": {
      "type": "object",
      "patternProperties": {
        ".+": {
          "type": "string"
        }
      },
      "additionalProperties": false
    },
    "status": {
      "type": "integer"
    }
  },
  "additionalProperties": false
}`
		}

		buff := []byte(resp)

		if len(buff) > 0 {
			buff = buff[:len(buff)]
		}

		if _, err := w.Write(buff); err != nil {
			w.WriteHeader(500)
		}
	}))

	c.URL = api.URL
}