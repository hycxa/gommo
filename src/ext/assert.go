package ext

import (
	"errors"
	"testing"
)

func Assert(condition bool) {
	if !condition {
		panic(errors.New(""))
	}
}

func AssertE(err error) {
	if err != nil {
		panic(err)
	}
}

func AssertT(t *testing.T, condition bool) {
	if !condition {
		t.Fatalf(Stack())
	}
}

func AssertB(b *testing.B, condition bool) {
	if !condition {
		b.Fatalf(Stack())
	}
}

func TestingAssert(t *testing.T, condition bool, err error) {
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

func (e MyError) Error() string {
	return e.ErrorStr
}
