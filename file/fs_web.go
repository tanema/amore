// +build js

package file

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type (
	amoreFS struct {
		zipFiles map[string]*zip.File
	}
	//create a struct so write can be implemented so that it statisfies ReadWriterCloser
	file struct {
		io.ReadCloser
		data   []byte // non-nil if regular file
		reader *io.SectionReader
		once   sync.Once
	}
	File interface {
		io.Reader
		io.Writer
		io.Seeker
		io.Closer
	}
)

var (
	zipData string
	fs      = &amoreFS{
		zipFiles: make(map[string]*zip.File),
	}
)

func Register(data string) {
	fs.register(data)
}

func (fs *amoreFS) register(data string) {
	zipReader, err := zip.NewReader(strings.NewReader(data), int64(len(data)))
	if err != nil {
		panic(err)
	}
	for _, file := range zipReader.File {
		fs.zipFiles[file.Name] = file
	}
}

func (fs *amoreFS) open(path string) (File, error) {
	path = fs.path(path)
	zipFile, ok := fs.zipFiles[path]
	if !ok {
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return &file{
			ReadCloser: resp.Body,
			data:       all,
			reader:     io.NewSectionReader(bytes.NewReader(all), 0, int64(len(all))),
		}, nil
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

func (fs *amoreFS) stat(path string) (os.FileInfo, error) {
	path = fs.path(path)
	zipFile, ok := fs.zipFiles[path]
	if !ok {
		return nil, nil
	}
	return zipFile.FileInfo(), nil
}

func (fs *amoreFS) readFile(path string) ([]byte, error) {
	path = fs.path(path)
	zipFile, ok := fs.zipFiles[path]
	if !ok {
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(resp.Body)
	}

	rc, err := zipFile.Open()
	defer rc.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(rc)
}

func (fs *amoreFS) mkDir(path string) error {
	return fmt.Errorf("cannot mkdir on web platform")
}

func (fs *amoreFS) remove(path string) error {
	return fmt.Errorf("cannot rm on web platform")
}

func (fs *amoreFS) ext(filename string) string {
	return filepath.Ext(fs.path(filename))
}

func (fs *amoreFS) path(filename string) string {
	p := strings.Replace(filename, "//", "/", -1)

	if _, ok := fs.zipFiles["assets/"+p]; ok {
		return "assets/" + p
	}

	return p
}

func (f *file) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("cannot write to file on web platform")
}

// Reads bytes into p, returns the number of read bytes.
func (f *file) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}

// Seeks to the offset.
func (f *file) Seek(offset int64, whence int) (ret int64, err error) {
	return f.reader.Seek(offset, whence)
}

// Returns an empty slice of files, directory
// listing is disabled.
func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	// directory listing is disabled.
	return make([]os.FileInfo, 0), nil
}
