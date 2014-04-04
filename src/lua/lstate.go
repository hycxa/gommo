package lua

/*
#cgo CFLAGS: -I/opt/local/include/
#cgo LDFLAGS: -llua -L/opt/local/lib
#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>
*/
import "C"

import ()

/*
#define LUA_TNONE		(-1)

#define LUA_TNIL		0
#define LUA_TBOOLEAN		1
#define LUA_TLIGHTUSERDATA	2
#define LUA_TNUMBER		3
#define LUA_TSTRING		4
#define LUA_TTABLE		5
#define LUA_TFUNCTION		6
#define LUA_TUSERDATA		7
#define LUA_TTHREAD		8
*/

type L struct {
	s *C.lua_State
}

func NewLua() *L {
	l := new(L)
	l.s = C.luaL_newstate()
	C.luaL_openlibs(l.s)
	return l
}

func (l *L) Close() {
	C.lua_close(l.s)
}

func value(s *C.lua_State, i C.int) interface{} {
	t := C.lua_type(s, i)
	switch t {
	case C.LUA_TNIL:
		return nil
	case C.LUA_TBOOLEAN:
		return bool(C.lua_toboolean(s, i) != 0)
	case C.LUA_TLIGHTUSERDATA:
		return nil
	case C.LUA_TNUMBER:
		return int64(C.lua_tonumberx(s, i, nil))
	case C.LUA_TSTRING:
		return C.GoString(C.lua_tolstring(s, i, nil))
	case C.LUA_TTABLE:
		return nil
	case C.LUA_TFUNCTION:
		return nil
	case C.LUA_TUSERDATA:
		return nil
	case C.LUA_TTHREAD:
		return nil
	}
	return nil
}

func (l *L) getRetValue(args C.int) *[]interface{} {
	ret := make([]interface{}, int(args))

	for i, index := C.int(1), 0; i <= args; i++ {
		ret[index] = value(l.s, i)
		index++
	}
	return &ret
}

func (l *L) DoString(str string) (ok bool, ret *[]interface{}) {
	n := C.lua_gettop(l.s)
	C.luaL_loadstring(l.s, C.CString(str))

	if C.lua_pcallk(l.s, 0, C.LUA_MULTRET, 0, 0, nil) == 0 {
		ok = true
	} else {
		ok = false
	}
	retCnt := C.lua_gettop(l.s) - n
	ret = l.getRetValue(retCnt)
	C.lua_settop(l.s, n)
	return
}

func (l *L) Call(f string, args ...interface{}) (ok bool, ret *[]interface{}) {
	n := C.lua_gettop(l.s)
	C.lua_getglobal(l.s, C.CString(f))
	nargs := 0
	for _, v := range args {
		switch v.(type) {
		case int:
			C.lua_pushinteger(l.s, C.lua_Integer(v.(int)))
		case int64:
			C.lua_pushinteger(l.s, C.lua_Integer(v.(int64)))
		case string:
			C.lua_pushlstring(l.s, C.CString(v.(string)), C.size_t(len(v.(string))))
		case bool:
			if v.(bool) {
				C.lua_pushboolean(l.s, C.int(1))
			} else {
				C.lua_pushboolean(l.s, C.int(0))
			}
		default:
			C.lua_pushnil(l.s)
		}
		nargs++
	}

	if C.lua_pcallk(l.s, C.int(nargs), C.LUA_MULTRET, 0, 0, nil) == C.LUA_OK {
		ok = true
	} else {
		ok = false
	}
	retCnt := C.lua_gettop(l.s) - n
	ret = l.getRetValue(retCnt)
	C.lua_settop(l.s, n)
	return
}
