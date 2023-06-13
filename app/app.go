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

// Generate - Return the command line application to be executed
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
			Action: setupToSearch,
		},
	}

	return app
}

func setupToSearch(c *cli.Context) {
	pathToSearch := c.String("path")
	content := c.String("contain")

	if len(pathToSearch) > 0 {
		if len(content) > 0 {
			printFiles := searching(pathToSearch, content)
			fmt.Println(printFiles)
		} else {
			log.Fatal("Try again! Specify the parameter contain!")
		}
	} else {
		log.Fatal("Try again! Specify the parameter path!")
	}
}

func searching(path, content string) []string {
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
	
	channelFiles := make(chan os.DirEntry, len(files))
	channelResults := make(chan string, len(files))

	for i := 0; i < runtime.NumCPU(); i++ {
		go searchingRoutine(path, content, reg, channelFiles, channelResults)
	}

	for _, file := range files {
		channelFiles <- file
	}
	close(channelFiles)

	var foundFiles []string
	for i := 0; i < len(files); i++ {
		found := <-channelResults
		if len(found) > 0 {
			foundFiles = append(foundFiles, found)
		}
	}

	return foundFiles
}

func searchingRoutine(path, content string, reg *regexp.Regexp, channelFiles <-chan os.DirEntry, channelResults chan<- string) {
	for file := range channelFiles {
		channelResults <- checkingFiles(file, path, content, reg)
	}
}

func checkingFiles(file os.DirEntry, path, content string, reg *regexp.Regexp) string {
	if reg.MatchString(file.Name()) {
		return file.Name()
	} else {
		if extension := strings.Split(file.Name(), "."); extension[len(extension)-1] == "txt" {
			data, err := os.ReadFile(path + "/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			foundContent := strings.Split(string(data), " ")
			for _, checking := range foundContent {
				if checking == content {
					return file.Name()
				}
			}
		}
	}
	return ""
}
