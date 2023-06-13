package main

import (
	"file-crawler/app"
	"log"
	"os"
)

func main() {
	app := app.Generate()
	if erro := app.Run(os.Args); erro != nil {
		log.Fatal(erro)
	}

}