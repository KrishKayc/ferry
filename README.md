

# jira - search
CLI utility to search and download issues (stories/bugs/tasks) From JIRA with a configurable 'filters' and return 'fields'. Retrieved issues from JIRA will be downloaded to 'output.csv'

## Command-line ##
**Usage**

    jira search [flags]

**Available Commands**
```
    search      Search and export Issues From JIRA
    help        Help about any command
```

**Flags**
```
  -h, --help   help for jira
```

Use "jira [command] --help" for more information about a command.


**Search**

**Simple Search**

```
 ./jira search --url "yoursite.atlassian.com" --project "Your Project" --issuetype Story 
```

**Filter multiple 'issue types'**

```
 ./jira search --url "yoursite.atlassian.com" --project "Your Project" --issuetype "Story,Bug,SubTask"
```
    
**Return specified fields**
```
 ./jira search --url "yoursite.atlassian.com" --project "Your Project" --issuetype "Story" --fields "summary,assignee,reporter,sprint" 
```

**Apply Filters**
*The below command will pull all the **Bugs** from **sprint5** assigned to the user **John***
```
 ./jira search --url "yoursite.atlassian.com" --project "Your Project" --issuetype "Bug" --fields "summary,sprint" --filters "sprint:sprint5,assignee:John"
```






