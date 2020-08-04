package main

import "goJIRA/jirafinder"

func main() {
	f := jirafinder.NewJiraFinder("config.json")
	f.Search()

}
