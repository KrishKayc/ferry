package main

import (
	"fmt"
	"github.com/gojira/ferry/jirafinder"
	"log"
)

func main() {
	err, f := jirafinder.NewJiraFinder("config.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := f.Search(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(" Download complete!!. Results exported to " + "'" + f.Config.DownloadPath + "'")
}
