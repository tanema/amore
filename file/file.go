// Package file is meant to take care of all asset file opening so that file
// access is safe if trying to access a bundled file or a file from the disk
package file

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/goxjs/glfw"
)

var (
	// map of zip file data related to thier real file path for consistent access
	zipFiles = make(map[string]*zip.File)
)

// Register will be called by bundled assets to register the bundled files into the
// zip file data to be used by the program.
func Register(data string) {
	zipReader, err := zip.NewReader(strings.NewReader(data), int64(len(data)))
	if err != nil {
		panic(err)
	}
	for _, file := range zipReader.File {
		zipFiles[file.Name] = file
	}
}

// Read will read a file at the path specified in total and return a byte
// array of the file contents
func Read(path string) ([]byte, error) {
	path = normalizePath(path)
	var file io.ReadCloser
	var err error
	if zipfile, ok := zipFiles[path]; ok {
		file, err = zipfile.Open()
	} else {
		file, err = Open(path)
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

// ReadString acts like Read but instead return a string. This is useful in certain
// circumstances.
func ReadString(filename string) string {
	s, err := Read(filename)
	if err != nil {
		return ""
	}
	return string(s[:])
}

// Open will return the file if its bundled or on disk and return a File interface
// for the file and an error if it does not exist. The File interface allows for
// consitent access to disk files and zip files.
func Open(path string) (io.ReadCloser, error) {
	path = normalizePath(path)
	zipFile, ok := zipFiles[path]
	if !ok {
		return glfw.Open(path)
	}
	return zipFile.Open()
}

// Ext will return the extention of the file
func Ext(filename string) string {
	return filepath.Ext(normalizePath(filename))
}

// normalizePath will prefix a path with assets/ if it comes from that directory
// or if it is bundled in such a way.
func normalizePath(filename string) string {
	p := strings.Replace(filename, "//", "/", -1)
	//bundle assets from asset folder
	if _, ok := zipFiles["assets/"+p]; ok {
		return "assets/" + p
	}
	return p
}
