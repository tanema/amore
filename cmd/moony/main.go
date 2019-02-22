package main

import (
	"log"
	"os"

	"github.com/tanema/amore/cmd"
)

func main() {
	exitStatus, err := cmd.App.Run()
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitStatus)
}
