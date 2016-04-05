package main

import (
	"fmt"
	"os"
	"path"
)

var assetFolders = []string{"images", "audio", "fonts", "shaders", "data"}

func newProject(project_name string) {
	for _, folder := range assetFolders {
		if err := os.MkdirAll(path.Join("assets", folder), os.ModeDir|os.ModePerm); err != nil {
			exitWithError(err)
		}
		fmt.Println(fmt.Sprintf("created assets/%v", folder))
	}
	err := writeOutTemplate("./conf.toml", configTemplate, struct{ Name string }{Name: project_name})
	if err != nil {
		exitWithError(err)
	}
	fmt.Println("created conf.toml")
	err = writeOutTemplate("./main.go", mainTemplate, struct{ PackageName string }{PackageName: *namePackage})
	if err != nil {
		exitWithError(err)
	}
	fmt.Println("created main.go")
}
