package main

import "jiraSearch_git/jirafinder"

func main() {
	//Read config from the 'config.json' file
	jirafinder.SetConfig("config.json")
	f := &jirafinder.JiraFinder{}

	f.Search()

}
