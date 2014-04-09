package lua

/*
#cgo CFLAGS: -I/opt/local/include/
#cgo LDFLAGS: -llua -L/opt/local/lib
#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>

int openlibs(lua_State* L) {
 	luaL_openlibs(L);
 	return 0;
}
*/
import "C"

import (
	"ext"
)

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

import "reflect"

type L struct {
	s *C.lua_State
}

func NewLua() *L {
	l := new(L)
	l.s = C.luaL_newstate()
	C.lua_pushcclosure(l.s, C.lua_CFunction(C.openlibs), 0)
	ext.Assert(C.lua_pcallk(l.s, 0, C.LUA_MULTRET, 0, 0, nil) == 0, "luaL_openlibs failed, not enough memory.")
	return l
}

func (l *L) Close() {
	C.lua_close(l.s)
}

func parseTable(s *C.lua_State, i C.int, retMap map[interface{}]interface{}) {
	C.lua_pushvalue(s, i)
	C.lua_pushnil(s)

	keyIndex := C.int(-2)
	valueIndex := C.int(-1)

	for int(C.lua_next(s, keyIndex)) != 0 {
		var key interface{}

		kt := C.lua_type(s, keyIndex)
		switch kt {
		case C.LUA_TNUMBER:
			key = int64(C.lua_tonumberx(s, keyIndex, nil))
		case C.LUA_TSTRING:
			key = C.GoString(C.lua_tolstring(s, keyIndex, nil))
		default:
			return
		}

		retMap[key] = value(s, valueIndex)

		C.lua_settop(s, -2)
	}
	C.lua_settop(s, -2)
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
		retMap := make(map[interface{}]interface{})
		parseTable(s, i, retMap)
		return retMap
	case C.LUA_TFUNCTION:
		return nil
	case C.LUA_TUSERDATA:
		return nil
	case C.LUA_TTHREAD:
		return nil
	}
	return nil
}

func (l *L) getRetValue(args C.int) []interface{} {
	ret := make([]interface{}, int(args))

	for i, index := C.int(1), 0; i <= args; i++ {
		ret[index] = value(l.s, i)
		index++
	}
	return ret
}

func (l *L) DoString(str string) (ok bool, ret []interface{}) {
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

func pushTable(l *L, v reflect.Value) {
	C.lua_createtable(l.s, 0, 0)
	for i := 0; i < v.Len(); i++ {
		pushValue(l, i+1)
		pushValue(l, v.Index(i).Interface())
		C.lua_settable(l.s, -3)
	}
}

func pushMap(l *L, v reflect.Value) {
	C.lua_createtable(l.s, 0, 0)
	keys := v.MapKeys()
	for i := 0; i < len(keys); i++ {
		pushValue(l, keys[i].Interface())
		pushValue(l, v.MapIndex(keys[i]).Interface())
		C.lua_settable(l.s, -3)
	}
}

func pushValue(l *L, v interface{}) {
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
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice:
			pushTable(l, rv)
		case reflect.Map:
			pushMap(l, rv)
		default:
			C.lua_pushnil(l.s)
		}
	}
}

func (l *L) Call(f string, args ...interface{}) (ok bool, ret []interface{}) {
	n := C.lua_gettop(l.s)
	C.lua_getglobal(l.s, C.CString(f))

	nargs := C.int(len(args))
	ext.Assert(C.lua_checkstack(l.s, nargs) != 0, "not enough free stack slots")
	for _, v := range args {
		pushValue(l, v)
	}

	if C.lua_pcallk(l.s, nargs, C.LUA_MULTRET, 0, 0, nil) == C.LUA_OK {
		ok = true
	} else {
		ok = false
	}
	retCnt := C.lua_gettop(l.s) - n
	ret = l.getRetValue(retCnt)
	C.lua_settop(l.s, n)
	return
}

