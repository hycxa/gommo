package lua

import (
	"ext"
	"testing"
)

func TestNewState(t *testing.T) {
	l := NewLua()
	l.Close()
}

func TestDoString(t *testing.T) {
	l := NewLua()
	defer l.Close()

	ok, r := l.DoString("function echo(...) return ... end")
	ext.AssertT(t, ok && len(*r) == 0, "dostring error")

	ok, r = l.DoString("return echo(1, 2, 3)")
	ok, r = l.DoString("return echo(1, \"s\", true)")

	ext.AssertT(t, ok && len(*r) == 3, "call error")
	ext.AssertT(t, 1 == (*r)[0].(int64), "return 1 error")
	ext.AssertT(t, "s" == (*r)[1].(string), "return 2 error")
	ext.AssertT(t, true == (*r)[2].(bool), "return 3 error")
}

func TestCall(t *testing.T) {
	l := NewLua()
	defer l.Close()

	ok, r := l.DoString("function echo(...) return ... end")
	ext.AssertT(t, ok && len(*r) == 0, "dostring error")

	ok, r = l.Call("echo", 4, "def ghit", true, "quit")
	ext.AssertT(t, ok && len(*r) == 4, "call error")
	ext.AssertT(t, 4 == (*r)[0].(int64), "return 1 error")
	ext.AssertT(t, "def ghit" == (*r)[1].(string), "return 2 error")
	ext.AssertT(t, true == (*r)[2].(bool), "return 3 error")
	ext.AssertT(t, "quit" == (*r)[3].(string), "return 4 error")
}



