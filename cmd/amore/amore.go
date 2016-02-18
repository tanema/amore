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
	namePackage    = "main"
	nameSourceFile = "asset_bundle.go"
	amoreVersion   = "0.0.1"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		if args[0] == "bundle" {
			bundle()
		} else if args[0] == "init" || args[0] == "new" {
			newProject()
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
