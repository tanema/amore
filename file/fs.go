package file

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// map of zip file data related to thier real file path for consistent access
	zipFiles = make(map[string]*zip.File)
	// map of paths to asset files for easy and fast path normilization
	assetFiles = make(map[string]string)
)

// init will walk the assets directory and take inventory to make normalizing
// asset paths a lot easier and faster
func init() {
	filepath.Walk("assets", func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() || strings.HasPrefix(fi.Name(), ".") {
			return nil
		}
		simple_path := strings.Replace(path, "assets/", "", -1)
		assetFiles[simple_path] = path
		return nil
	})
}

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

// stat will return the file info for the path provided or an error if there was
// an error retrieving it.
func stat(path string) (os.FileInfo, error) {
	path = normalizePath(path)
	zipFile, ok := zipFiles[path]
	if !ok {
		return os.Stat(path)
	}
	return zipFile.FileInfo(), nil
}

// normalizePath will prefix a path with assets/ if it comes from that directory
// or if it is bundled in such a way.
func normalizePath(filename string) string {
	p := strings.Replace(filename, "//", "/", -1)

	//bundle assets from asset folder
	if _, ok := zipFiles["assets/"+p]; ok {
		return "assets/" + p
	}

	//asset helper so user doesnt need to prefix assets with  assets/ path
	if _, ok := assetFiles[p]; ok {
		return assetFiles[p]
	}

	return p
}

// file statisfies the ReadWriterCloser interface and can be use as File and imitate
// os.File
type file struct {
	io.ReadCloser
	data   []byte // non-nil if regular file
	reader *io.SectionReader
	once   sync.Once
}

// Write is a stub method because you cannot change a bundled file, it will always
// return an error say as such.
func (f *file) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("Cannot write to a bundled file")
}

// Read reads bytes into p and returns the number of read bytes and an error if
// there was one while reading.
func (f *file) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}

// Seeks, calls seek on the data reader to seek to the offset in the data from the zip file.
func (f *file) Seek(offset int64, whence int) (ret int64, err error) {
	return f.reader.Seek(offset, whence)
}

// Returns an empty slice of files, directory
// listing is disabled for bundled files.
func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	// directory listing is disabled.
	return make([]os.FileInfo, 0), nil
}
