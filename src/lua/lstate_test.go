package lua

import (
	"ext"
	"fmt"
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

	ok, r := l.DoString(`function echo(...)
		return ... 
	end`)
	ext.AssertT(t, ok && len(r) == 0, "dostring error")

	tab := make([]int, 5)
	tab[0] = 5
	tab[1] = 3
	tab[2] = 2
	tab[3] = 1
	tab[4] = 4

	mmap := make(map[int]int)
	mmap[1] = 22
	mmap[8] = 7
	mmap[9] = 16
	mmap[5] = 33

	ok, r = l.Call("echo", 4, "def ghit", true, tab, "quit", mmap)
	ext.AssertT(t, ok && len(r) == 6, "call error")
	ext.AssertT(t, 4 == r[0].(int64), "return 1 error")
	ext.AssertT(t, "def ghit" == r[1].(string), "return 2 error")
	ext.AssertT(t, true == r[2].(bool), "return 3 error")
	ext.AssertT(t, "quit" == r[4].(string), "return 5 error")
	fmt.Println(r[3])
	fmt.Println(r[5])

	retTab := r[3].(map[int]int)
	for i := 0; i < len(tab); i++ {
		ext.AssertT(t, tab[i] == retTab[i], "return tab error")
	}
}

func BenchmarkCallEffInt(b *testing.B) {
	l := NewLua()
	defer l.Close()

	l.DoString("function echo(...) return ... end")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Call("echo", 4)
	}
}

func BenchmarkCallEffStr(b *testing.B) {
	l := NewLua()
	defer l.Close()

	l.DoString("function echo(...) return ... end")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Call("echo", "abc")
	}
}

func BenchmarkCallEffBool(b *testing.B) {
	l := NewLua()
	defer l.Close()

	l.DoString("function echo(...) return ... end")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Call("echo", false)
	}
}

func BenchmarkCallEffBaseComplex(b *testing.B) {
	l := NewLua()
	defer l.Close()

	l.DoString("function echo(...) return ... end")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Call("echo", 1, "testecho", true)
	}
}

func BenchmarkCallEffSlice(b *testing.B) {
	l := NewLua()
	defer l.Close()

	l.DoString("function echo(...) return ... end")
	tab := make([]int, 5)
	tab[0] = 5
	tab[1] = 3
	tab[2] = 2
	tab[3] = 1
	tab[4] = 4
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Call("echo", tab)
	}
}

func BenchmarkCallEffMap(b *testing.B) {
	l := NewLua()
	defer l.Close()

	l.DoString("function echo(...) return ... end")
	mmap := make(map[int]int)
	mmap[1] = 22
	mmap[8] = 7
	mmap[9] = 16
	mmap[5] = 33

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Call("echo", mmap)
	}
}
