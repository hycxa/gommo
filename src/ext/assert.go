package ext

import (
	"fmt"
	"testing"
)

func Assert(condition bool) {
	if !condition {
		panic("assert failed")
	}
}

func AssertM(condition bool, f string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf(f, v...))
	}
}

func AssertE(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func AssertT(t *testing.T, condition bool, f string, v ...interface{}) {
	if !condition {
		s := fmt.Sprintf(f, v...)
		t.Fatalf("%s\n%s", s, Stack())
	}
}
