/*
 * This package is meant to take care of all asset file opening
 * In the future this will be made to use a pre configured directory
 * also when put onto mobile it will be made to access thos appropriately as well
 */
package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func Read(filename string) ([]byte, error) {
	buf, err := ioutil.ReadFile(filterPath(filename))
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func ReadString(filename string) string {
	s, err := Read(filename)
	if err != nil {
		panic(err)
	}
	return string(s[:])
}

func NewFileData(filename string) (os.FileInfo, error) {
	info, err := os.Stat(filterPath(filename))
	if err != nil {
		return nil, err
	}
	return info, nil
}

func NewFile(filename string) (*os.File, error) {
	f, err := os.Open(filterPath(filename))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func CreateDirectory(filename string) error {
	return os.MkdirAll(filterPath(filename), os.ModeDir)
}

func Exists(filename string) bool {
	info, err := os.Stat(filterPath(filename))
	return info != nil && err == nil
}

func Remove(filename string) error {
	err := os.Remove(filterPath(filename))
	if err != nil {
		return err
	}
	return nil
}

func IsDirectory(filename string) bool {
	info, err := os.Stat(filterPath(filename))
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsFile(filename string) bool {
	info, err := os.Stat(filterPath(filename))
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func IsSymLink(filename string) bool {
	info, err := os.Stat(filterPath(filename))
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeSymlink) != 0
}

func GetSize(filename string) int32 {
	info, err := os.Stat(filterPath(filename))
	if err != nil {
		return 0
	}
	return int32(info.Size())
}

func GetLastModified(filename string) time.Time {
	info, err := os.Stat(filterPath(filename))
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

func Ext(filename string) string {
	return filepath.Ext(filterPath(filename))
}

func filterPath(filename string) string {
	if !filepath.IsAbs(filename) {
		filename = filepath.Join("assets", filename)
	}
	return filename
}
