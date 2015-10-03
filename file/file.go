package file

import (
	"io/ioutil"
)

func Read(filename string) string {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(buf[:])
}
