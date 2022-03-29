//go:build !solution
// +build !solution

package fileleak

import (
	"io/ioutil"
	"log"
	"reflect"
)

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func getOpenedFileNames() []string {
	files, err := ioutil.ReadDir("/proc/self/fd")
	if err != nil {
		log.Fatal(err)
	}

	names := make([]string, 0, len(files))
	for _, file := range files {
		names = append(names, file.Name())
	}
	return names
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
