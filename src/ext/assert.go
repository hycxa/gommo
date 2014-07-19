package ext

import (
	"fmt"
	"testing"
)

func Assert(condition bool, f string, v ...interface{}) {
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

func TestingAssert(t* testing.T, condition bool, err error) {
	if !condition {
		if err != nil {
			t.Fatalf("%s\n%s", err.Error(), Stack())
		} else {
			t.Fatalf(Stack())
		}
	}
}

type MyError struct {
	ErrorStr string
}

func (e MyError) Error() string{
	return e.ErrorStr
}
