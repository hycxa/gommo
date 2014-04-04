package lua

import (
	"ext"
	"testing"
)

func TestNewState(t *testing.T) {
	l := NewL()
	l.Close()
}

func TestCall(t *testing.T) {
	l := NewL()
	defer l.Close()

	r := l.DoString("function echo(...) return 1 , 2, 3 end")
	ext.AssertT(t, len(r) == 0, "dostring error")

	r = l.DoString("echo(1, 2, 3)")

	ext.AssertT(t, len(r) == 1, "echo error")
}
