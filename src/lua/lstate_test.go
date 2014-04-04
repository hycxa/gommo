package lua

import (
	"ext"
	"testing"
)

func TestNewState(t *testing.T) {
	l := NewL()
	l.Close()
}

func TestDoString(t *testing.T) {
	l := NewL()
	defer l.Close()

	r := l.DoString("function echo(...) return 1 , 2, 3 end")
	ext.AssertT(t, len(r) == 0, "dostring error")

	r = l.DoString("echo(1, 2, 3)")

	ext.AssertT(t, len(r) == 1, "echo error")
}

func TestCall(t *testing.T) {
	l := NewL()
	defer l.Close()

	r := l.DoString("function echo(...) return ... end")
	ext.AssertT(t, len(r) == 0, "dostring error")

	r = l.Call("echo", 1, "s", true)
	ext.AssertT(t, len(r) == 3, "call error")
	ext.AssertT(t, 1 == r[0].(int), "return 1 error")
	ext.AssertT(t, "s" == r[1].(string), "return 2 error")
	ext.AssertT(t, true == r[2].(bool), "return 3 error")
}
