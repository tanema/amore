package main

import (
	"fmt"
	"os"
	"path"
)

var assetFolders = []string{"images", "audio", "fonts", "shaders", "data"}

// newProject will generate a file structure for a new project, including a main.go
// file, a default config and all the asset directories.
func newProject(project_name string) {
	for _, folder := range assetFolders {
		if err := os.MkdirAll(path.Join("assets", folder), os.ModeDir|os.ModePerm); err != nil {
			exitWithError(err)
		}
		fmt.Println(fmt.Sprintf("generated -> assets/%v directory", folder))
	}
	err := writeOutTemplate("./conf.toml", configTemplate, struct{ Name string }{Name: project_name})
	if err != nil {
		exitWithError(err)
	}
	fmt.Println("generated -> conf.toml")
	err = writeOutTemplate("./main.go", mainTemplate, struct{ PackageName string }{PackageName: *namePackage})
	if err != nil {
		exitWithError(err)
	}
	fmt.Println("generated -> main.go")
}
