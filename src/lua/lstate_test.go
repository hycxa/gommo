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

	ok, r := l.DoString("function echo(...) return ... end")
	ext.AssertT(t, ok && len(r) == 0, "dostring error")

	ok, r = l.DoString("return echo(1, \"s\", true)")

	ext.AssertT(t, ok && len(r) == 3, "call error")
	ext.AssertT(t, 1 == r[0].(int64), "return 1 error")
	ext.AssertT(t, "s" == r[1].(string), "return 2 error")
	ext.AssertT(t, true == r[2].(bool), "return 3 error")
}

func TestCall(t *testing.T) {
	l := NewL()
	defer l.Close()

	ok, r := l.DoString("function echo(...) return ... end")
	ext.AssertT(t, len(r) == 0, "dostring error")

	ok, r = l.Call("echo", 1, "s", true)
	ext.AssertT(t, ok && len(r) == 3, "call error")
	ext.AssertT(t, 1 == r[0].(int64), "return 1 error")
	ext.AssertT(t, "s" == r[1].(string), "return 2 error")
	ext.AssertT(t, true == r[2].(bool), "return 3 error")
}
