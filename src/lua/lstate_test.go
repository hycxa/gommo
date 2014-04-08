package lua

import (
	"ext"
	"testing"
)

func TestNewLua(t *testing.T) {
	l := NewLua()
	l.Close()
}

func BenchmarkNewLua(b *testing.B) {
	l := make([]*L, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l[i] = NewLua()
		//defer l[i].Close()
	}
}

func TestDoString(t *testing.T) {
	l := NewLua()
	defer l.Close()

	ok, r := l.DoString("function echo(...) return ... end")
	ext.AssertT(t, ok && len(r) == 0, "dostring error")

	ok, r = l.DoString("return echo(1, 2, 3)")
	ok, r = l.DoString("return echo(1, \"s\", true)")

	ext.AssertT(t, ok && len(r) == 3, "call error")
	ext.AssertT(t, 1 == r[0].(int64), "return 1 error")
	ext.AssertT(t, "s" == r[1].(string), "return 2 error")
	ext.AssertT(t, true == r[2].(bool), "return 3 error")
}

func TestCall(t *testing.T) {
	l := NewLua()
	defer l.Close()

	ok, r := l.DoString("function echo(...) return ... end")
	ext.AssertT(t, ok && len(r) == 0, "dostring error")

	ok, r = l.Call("echo", 4, "def ghit", true, "quit")
	ext.AssertT(t, ok && len(r) == 4, "call error")
	ext.AssertT(t, 4 == r[0].(int64), "return 1 error")
	ext.AssertT(t, "def ghit" == r[1].(string), "return 2 error")
	ext.AssertT(t, true == r[2].(bool), "return 3 error")
	ext.AssertT(t, "quit" == r[3].(string), "return 4 error")
}

func BenchmarkCallEff1(b *testing.B) {
	l := NewLua()
	defer l.Close()

	l.DoString("function echo(...) return ... end")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Call("echo", 4)
	}
}

func BenchmarkCallEff2(b *testing.B) {
	l := NewLua()
	defer l.Close()

	l.DoString("function echo(...) return ... end")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Call("echo", "hello")
	}
}

func BenchmarkCallEff3(b *testing.B) {
	l := NewLua()
	defer l.Close()

	l.DoString("function echo(...) return ... end")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Call("echo", false)
	}
}

func BenchmarkCallEff4(b *testing.B) {
	l := NewLua()
	defer l.Close()

	l.DoString("function echo(...) return ... end")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Call("echo", 1, "testecho", true)
	}
}
