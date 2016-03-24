/*
The file Package is meant to take care of all asset file opening so that file
access is safe if trying to access a bundled file or a file from the disk
*/
package file

import (
	"os"
	"time"
)

func Read(filename string) ([]byte, error) {
	return fs.readFile(filename)
}

func ReadString(filename string) string {
	s, err := Read(filename)
	if err != nil {
		panic(err)
	}
	return string(s[:])
}

func NewFileData(filename string) (os.FileInfo, error) {
	return fs.stat(filename)
}

func NewFile(filename string) (File, error) {
	return fs.open(filename)
}

func CreateDirectory(filename string) error {
	return fs.mkDir(filename)
}

func Exists(filename string) bool {
	info, err := fs.stat(filename)
	return info != nil && err == nil
}

func Remove(filename string) error {
	return fs.remove(filename)
}

func IsDirectory(filename string) bool {
	info, err := fs.stat(filename)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsFile(filename string) bool {
	return !IsDirectory(filename)
}

func IsSymLink(filename string) bool {
	info, err := fs.stat(filename)
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeSymlink) != 0
}

func GetSize(filename string) int32 {
	info, err := fs.stat(filename)
	if err != nil {
		return 0
	}
	return int32(info.Size())
}

func GetLastModified(filename string) time.Time {
	info, err := fs.stat(filename)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

func Ext(filename string) string {
	return fs.ext(filename)
}
