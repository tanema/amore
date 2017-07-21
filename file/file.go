/*
The file Package is meant to take care of all asset file opening so that file
access is safe if trying to access a bundled file or a file from the disk
*/
package file

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// the File interface makes accessing bundled files and os.File consistent
type File interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

// Read will read a file at the path specified in total and return a byte
// array of the file contents
func Read(path string) ([]byte, error) {
	path = normalizePath(path)
	zipFile, ok := zipFiles[path]
	if !ok {
		return ioutil.ReadFile(path)
	}

	rc, err := zipFile.Open()
	defer rc.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(rc)
}

// ReadString acts like Read but instead return a string. This is useful in certain
// circumstances.
func ReadString(filename string) string {
	s, err := Read(filename)
	if err != nil {
		panic(err)
	}
	return string(s[:])
}

// NewFileData retreives the file info for the path provided. If the file does
// not exist it will return an error
func NewFileData(filename string) (os.FileInfo, error) {
	return stat(filename)
}

// NewFile will return the file if its bundled or on disk and return a File interface
// for the file and an error if it does not exist. The File interface allows for
// consitent access to disk files and zip files.
func NewFile(path string) (File, error) {
	path = normalizePath(path)
	zipFile, ok := zipFiles[path]
	if !ok {
		return os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0777)
	}

	rc, err := zipFile.Open()
	if err != nil {
		return nil, err
	}
	all, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return &file{
		ReadCloser: rc,
		data:       all,
		reader:     io.NewSectionReader(bytes.NewReader(all), 0, zipFile.FileInfo().Size()),
	}, nil
}

// Create will create and return a new empty file at the pass path
func Create(path string) (File, error) {
	path = normalizePath(path)
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
}

// CreateDirectory will create all directories in the path given if they do not exist
func CreateDirectory(path string) error {
	return os.MkdirAll(normalizePath(path), os.ModeDir|os.ModePerm)
}

// Exists will return true if the file exists at the path provided and false if
// the file does not exist.
func Exists(filename string) bool {
	info, err := stat(filename)
	return info != nil && err == nil
}

// Remove will delete a file at the given path and return an error if there was
// and issue
func Remove(path string) error {
	return os.Remove(normalizePath(path))
}

// IsDirectory will return true if the path provided is a directory and false
// if the file does not exist or is not a directory.
func IsDirectory(filename string) bool {
	info, err := stat(filename)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile will return false if file does not exist or is directory. It will return
// true otherwise.
func IsFile(filename string) bool {
	info, err := stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// IsSymLink will return false if the file does not exist or is not a symlink. It
// will return true if the file is a symlink
func IsSymLink(filename string) bool {
	info, err := stat(filename)
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeSymlink) != 0
}

// GetSize will return the files size in bytes, it will return 0 if the file does
// not exist
func GetSize(filename string) int32 {
	info, err := stat(filename)
	if err != nil {
		return 0
	}
	return int32(info.Size())
}

// GetLastModified will return the time of when the file was last modified, if
// the file does not exist the time will be 0
func GetLastModified(filename string) time.Time {
	info, err := stat(filename)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

// Ext will return the extention of the file
func Ext(filename string) string {
	return filepath.Ext(normalizePath(filename))
}
