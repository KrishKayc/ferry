package main

import "github.com/gojira/ferry/jirafinder"

func main() {
	f := jirafinder.NewJiraFinder("config.json")
	f.Search()
}
