package main

import (
	"file-crawler/app"
	"log"
	"os"
)

func main() {
	generate := app.Generate()
	if err := generate.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
