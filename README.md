# ferry - JIRA search
Utility to Search and download Issues From JIRA with a configurable Filters and return fields. Retrieved Issues from JIRA will be downloaded as CSV in the path specified in the config.json file.

## Command-line ##
**Usage**

    ferry [flags]

    ferry [command]

**Available Commands**
```
    export      Search and export Issues From JIRA
    help        Help about any command
    version     Print the version
```

**Flags**
```
  -h, --help   help for ferry
```

Use "ferry [command] --help" for more information about a command.


**export command**
```
ferry export --config config.json --project "Your Project" --output ~/Documents/ferry.csv
```

**config.json** file specifies.

    * Filters to be applied. Example : Project, Issue Type, Sprint etc
    * FieldsToRetrive to be rendered as columns in the downloaded csv file

    

[![Build Status](https://travis-ci.org/KrishKayc/ferry.svg?branch=master)](https://travis-ci.org/KrishKayc/ferry)  [![codecov](https://codecov.io/gh/KrishKayc/ferry/branch/master/graph/badge.svg)](https://codecov.io/gh/KrishKayc/ferry)      [![Go Report Card](https://goreportcard.com/badge/github.com/KrishKayc/ferry)](https://goreportcard.com/report/github.com/KrishKayc/ferry)

## Output ##
**Download InProgress**:

![Output](https://github.com/KrishKayc/goJIRA/blob/master/output_screenshots/jiraSearch_finaloutput1.jpg)

**Download Complete**:

![FinalOutput](https://github.com/KrishKayc/goJIRA/blob/master/output_screenshots/jiraSearch_finaloutput2.jpg)

