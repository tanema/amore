// The bundle code is
// Licensed under the Apache License, Version 2.0 (the "License");
// as it came from github.com/rakyll/statik
package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	amoreVersion = "0.0.3"
)

var (
	namePackage    = flag.String("pkg", "main", "name of the go package for the source to be generated in")
	nameSourceFile = flag.String("out", "asset_bundle.go", "name of the outputted file for bundling")
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		if args[0] == "bundle" {
			bundle(args[1:]...)
		} else if args[0] == "init" || args[0] == "new" {
			project_name := ""
			if len(args) > 1 {
				project_name = args[1]
			}
			newProject(project_name)
		} else if args[0] == "version" {
			printVersion()
		}
	} else {
		printVersion()
	}
}

func printVersion() {
	fmt.Printf("Amore version: %v", amoreVersion)
}

// Prints out the error message and exists with a non-success signal.
func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
