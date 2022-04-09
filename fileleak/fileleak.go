//go:build !solution

package fileleak

import (
	"io/fs"
	"io/ioutil"
	"log"
	"reflect"
)

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func getOpenedFileNames() []fs.FileInfo {
	files, err := ioutil.ReadDir("/proc/self/fd")
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func VerifyNone(t testingT) {
	before := getOpenedFileNames()
	t.Cleanup(func() {
		after := getOpenedFileNames()
		if !reflect.DeepEqual(after, before) {
			t.Errorf("there are file leaks")
		}
	})
}
