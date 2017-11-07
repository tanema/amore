package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// bundler contains the bundled data and writes out the bundled data to the asset
// bundle
type bundler struct {
	buffer    *bytes.Buffer
	zipWriter *zip.Writer
}

// bundle will bundle inputed file paths. If no paths were given then it will bundle
// the default assets folder and the conf.toml file.
func bundle(inputs ...string) {
	b := newBundler()

	if inputs == nil || len(inputs) == 0 {
		inputs = []string{"assets", "conf.toml"}
	}

	for _, input := range inputs {
		fi, err := os.Stat(input)
		if err != nil {
			exitWithError(err)
		}
		if fi.IsDir() {
			err = b.addDir(input)
		} else {
			err = b.addFile(input)
		}
		if err != nil {
			exitWithError(err)
		}
	}

	err := b.writeOut()
	if err != nil {
		exitWithError(err)
	}
}

// newBundle creates a pointer to a new bundler object.
func newBundler() *bundler {
	nb := &bundler{
		buffer: new(bytes.Buffer),
	}
	nb.zipWriter = zip.NewWriter(nb.buffer)
	return nb
}

// addDir will walk the directory and add each individual file in the directory
// recursively.
func (bndlr *bundler) addDir(srcPath string) error {
	return filepath.Walk(srcPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Ignore directories and hidden files.
		// No entry is needed for directories in a zip file.
		// Each file is represented with a path, no directory
		// entities are required to build the hierarchy.
		if fi.IsDir() || strings.HasPrefix(fi.Name(), ".") {
			return nil
		}
		return bndlr.addFile(path)
	})
}

// addFile will write the file to the zipWriter to be written out later.
func (bndlr *bundler) addFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	f, err := bndlr.zipWriter.Create(path)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err == nil {
		fmt.Println(fmt.Sprintf("bundled file -> %v", path))
	}
	return err
}

// writeOut will write out the data to a asset_bundle file using a template
func (bndlr *bundler) writeOut() error {
	err := bndlr.zipWriter.Close()
	if err != nil {
		return err
	}

	zipData := FprintZipData(bndlr.buffer.Bytes())
	return writeOutTemplate(*nameSourceFile, bundleTemplate, struct {
		PackageName string
		Data        *bytes.Buffer
	}{
		PackageName: *namePackage,
		Data:        zipData,
	})
}

// FprintZipData Converts zip binary contents to a string literal.
func FprintZipData(zipData []byte) *bytes.Buffer {
	dest := new(bytes.Buffer)
	for _, b := range zipData {
		if b == '\n' {
			dest.WriteString(`\n`)
			continue
		}
		if b == '\\' {
			dest.WriteString(`\\`)
			continue
		}
		if b == '"' {
			dest.WriteString(`\"`)
			continue
		}
		if (b >= 32 && b <= 126) || b == '\t' {
			dest.WriteByte(b)
			continue
		}
		fmt.Fprintf(dest, "\\x%02x", b)
	}
	return dest
}
