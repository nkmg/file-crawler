package app

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"

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
		} else {
			log.Fatal("Try again! Specify the parameter contain!")
		}
	} else {
		log.Fatal("Try again! Specify the parameter path!")
	}
}

func searching(path, content string) {
	channel_files := make(chan os.DirEntry)
	channel_results := make(chan string)

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

	for i := 0; i < runtime.NumCPU(); i++ {
		go searching_routine(path, content, reg, channel_files, channel_results)
	}

	for _, file := range files {
		channel_files <- file
	}

	close(channel_files)

	for found := range channel_results {
		fmt.Println(found)
	}

}

func searching_routine(path, content string, reg *regexp.Regexp, channel_files <-chan os.DirEntry, channel_results chan<- string) {
	for file := range channel_files {
		if file_found := checking_files(file, path, content, reg); len(file_found) > 0 {
			channel_results <- file_found
		}
	}

}

func checking_files(file os.DirEntry, path, content string, reg *regexp.Regexp) string {
	if reg.MatchString(file.Name()) {
		return file.Name()
	} else {
		if extension := strings.Split(file.Name(), "."); extension[len(extension)-1] == "txt" {
			data, err := os.ReadFile(path + "/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			found_content := strings.Split(string(data), " ")
			for _, checking := range found_content {
				if checking == content {
					return file.Name()
				}
			}
		}
	}
	return ""
}
