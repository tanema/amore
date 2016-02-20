package main

import (
	"os"
	"path"
)

var assetFolders = []string{"images", "audio", "fonts", "shaders"}

func newProject(project_name string) {
	for _, folder := range assetFolders {
		if err := os.MkdirAll(path.Join("assets", folder), os.ModeDir); err != nil {
			exitWithError(err)
		}
	}
	err := writeOutTemplate("./conf.toml", configTemplate, struct{ Name string }{Name: project_name})
	if err != nil {
		exitWithError(err)
	}
	err = writeOutTemplate("./main.go", mainTemplate, struct{ PackageName string }{PackageName: *namePackage})
	if err != nil {
		exitWithError(err)
	}
}
