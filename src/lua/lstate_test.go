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
		l[i].Close()
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

	/*
		retTab := r[3].(map[int]int)
		for i := 0; i < len(tab); i++ {
			ext.AssertT(t, tab[i] == retTab[i], "return tab error")
		}
	*/
}

func TestCallCAndLua(t *testing.T) {
	l := NewLua()
	defer l.Close()
	l.InstallFunc()

	ok, r := l.DoString(`function echo()
		a = array.new(3, 5)
		assert(a:getx() == 3)
		assert(a:gety() == 5)
		a:setx(88)
		a:sety(99)
		assert(a:getx() == 88)
		assert(a:gety() == 99)
	end`)
	ext.AssertT(t, ok && len(r) == 0, "dostring error")

	ok, r = l.Call("echo")
	ext.AssertT(t, ok && len(r) == 0, "call error")
}

func BenchmarkCallEffCAndLua(b *testing.B) {
	l := NewLua()
	defer l.Close()
	l.InstallFunc()

	l.DoString(`
		a = array.new(3, 5)
		assert(a:getx() == 3)
		assert(a:gety() == 5)
	function echo()
		for i = 1, 1000 do
			a:getx()
		end
	end`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Call("echo")
	}
}

var testStr = `
local luatobin = require("lbin")
local t = {
	baggins = true,
	age = 24,
	name = 'dumbo' ,
	ponies = {'whisper','fartarse'},
	tt = 0,
	father = {
		baggins = false,
		age = 77,
		tt2 = 0,
		name = 'Wombo',
	}
}

function retEnc()
return luatobin.serialize(t, false)
end

function dec(s)
local to, to1 = luatobin.deserialize(s)
assert(to1.name == t.name)
assert(to1.father.age == t.father.age)
to1.tt = to1.tt + 1
return luatobin.serialize(to1, false)
end`

func TestLuaEnc(t *testing.T) {
	l := NewLua()
	defer l.Close()

	l.DoString(testStr)

	refEnc := l.GetRef("retEnc")
	ok, rs := l.CallRef(refEnc, nil)
	ext.AssertT(t, ok, "CallRef retEnc")

	refDec := l.GetRef("dec")
	ok, _ = l.CallRef(refDec, rs)
	ext.AssertT(t, ok, "CallRef refDec")
}

func BenchmarkLuaEncTab(b *testing.B) {
	l := NewLua()
	defer l.Close()

	l.DoString(testStr)
	ok, rs := l.Call("retEnc")

	ext.Assert(ok, "fffff")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, rs = l.Call("dec", rs)
	}
}

func BenchmarkLuaEnc(b *testing.B) {
	l := NewLua()
	defer l.Close()

	l.DoString(testStr)

	refEnc := l.GetRef("retEnc")
	ok, rs := l.CallRef(refEnc, nil)
	ext.Assert(ok, "ffff")

	refDec := l.GetRef("dec")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, rs = l.CallRef(refDec, rs)
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

	l.DoString(testStr)

	refEnc := l.GetRef("retEnc")
	_, rs := l.CallRef(refEnc, nil)
	s := string(rs)

	l.DoString(`function echo(str)
	return str
	end`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, rt := l.Call("echo", s)
		r := []rune(rt[0].(string))
		r[len(r)-1] = r[len(r)-1] + 1
		s = string(r)
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
