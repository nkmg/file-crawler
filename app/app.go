package app

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli"
)

// Return the command line application to be executed
func Generate() *cli.App {

	app := cli.NewApp()
	app.Name = "File Crawler"
	app.Usage = "Search for files with a specific filename or content name"
	app.Commands = []cli.Command{
		{
			Name:  "search",
			Usage: "./file-crawler <command> [arguments]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path",
					Usage: "Specific path to search. If not defined will search in all filesystem",
					Value: "",
				},
				cli.StringFlag{
					Name:  "contain",
					Usage: "Name that will searching in filenames and files",
					Value: "",
				},
			},
			Action: setup_to_search,
		},
	}

	return app
}

func setup_to_search(c *cli.Context) {
	path_to_search := c.String("path")
	content := c.String("contain")

	if len(path_to_search) > 0 {
		if len(content) > 0 {
			searching(path_to_search, content)
		}
	} else {
		log.Fatal("Try again specifying the path to search")
	}
}

func searching(path, content string) {
	channel := make(chan string)

	dir, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer dir.Close()

	files, err := dir.ReadDir(-1)
	if err != nil {
		log.Fatal(err)
	}

	reg, err := regexp.Compile(content)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for _, file := range files {
			if reg.MatchString(file.Name()) {
				channel <- file.Name()
				time.Sleep(time.Second * 1)
			}
		}
	}()

	go func() {
		for _, file := range files {
			if extension := strings.Split(file.Name(), "."); extension[len(extension)-1] == "txt" {
				data, err := os.ReadFile(path + "/" + file.Name())
				if err != nil {
					log.Fatal(err)
				}
				found_content := strings.Split(string(data), " ")
				for _, checking := range found_content {
					if checking == content {
						channel <- file.Name()
						time.Sleep(time.Second * 2)
					}
				}
			}
		}
		close(channel)
	}()

	for found := range channel {
		fmt.Println(found)
	}

}
