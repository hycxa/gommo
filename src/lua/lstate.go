package lua

/*
#cgo CFLAGS: -I /usr/include/lua5.1
#cgo LDFLAGS: -llua5.1
#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>
#include <stdlib.h>
#include "lua_compatible.h"

#include "userdata.h"
*/
import "C"

import (
	"ext"
	"reflect"
	"unsafe"
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

type L struct {
	s *C.lua_State
}

func NewLua() *L {
	l := new(L)
	l.s = C.luaL_newstate()
	C.lua_pushcclosure(l.s, C.lua_CFunction(C.openlibs), 0)
	ext.Assert(C.lpcall(l.s, 0, C.LUA_MULTRET, 0, 0, nil) == 0, "luaL_openlibs failed, not enough memory.")
	return l
}

func (l *L) InstallFunc() {
	C.install_func(l.s)
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
			key = int64(C.ltonumber(s, keyIndex))
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
		return int64(C.ltonumber(s, i))
	case C.LUA_TSTRING:
		slen := C.lrawlen(s, i)
		return C.GoStringN((C.lua_tolstring(s, i, &slen)), C.int(slen))
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
	Cs := C.CString(str)
	defer C.free(unsafe.Pointer(Cs))
	C.luaL_loadstring(l.s, Cs)

	if C.lpcall(l.s, 0, C.LUA_MULTRET, 0, 0, nil) == 0 {
		ok = true
	} else {
		ok = false
	}
	retCnt := C.lua_gettop(l.s) - n
	ret = l.getRetValue(retCnt)
	C.lua_settop(l.s, n)
	return
}

func (l *L) Loadfile(fileName string) (ok bool, ret []interface{}) {
	n := C.lua_gettop(l.s)
	Cs := C.CString(fileName)
	defer C.free(unsafe.Pointer(Cs))
	C.luaL_loadfile(l.s, Cs)

	if C.lpcall(l.s, 0, C.LUA_MULTRET, 0, 0, nil) == 0 {
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
		str := v.(string)
		Cs := C.CString(str)
		defer C.free(unsafe.Pointer(Cs))
		C.lua_pushlstring(l.s, Cs, C.size_t(len(str)))
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
	Cs := C.CString(f)
	defer C.free(unsafe.Pointer(Cs))
	C.lgetglobal(l.s, Cs)

	nargs := C.int(len(args))
	ext.Assert(C.lua_checkstack(l.s, nargs) != 0, "not enough free stack slots")
	for _, v := range args {
		pushValue(l, v)
	}

	if C.lpcall(l.s, nargs, C.LUA_MULTRET, 0, 0, nil) == C.LUA_OK {
		ok = true
	} else {
		ok = false
	}
	retCnt := C.lua_gettop(l.s) - n
	ret = l.getRetValue(retCnt)
	C.lua_settop(l.s, n)
	return
}

func (l *L) GetRef(s string) int {
	Cs := C.CString(s)
	defer C.free(unsafe.Pointer(Cs))
	//C.luaL_loadstring(l.s, Cs)
	C.lgetglobal(l.s, Cs)
	return int(C.luaL_ref(l.s, C.LUA_REGISTRYINDEX))
}

func isStackFunc(s *C.lua_State, index C.int) bool {
	return C.lua_type(s, index) == C.LUA_TFUNCTION
}

func (l *L) CallRef(funi int, args []byte) (ok bool, ret []byte) {
	n := C.lua_gettop(l.s)
	defer C.lua_settop(l.s, n)
	C.lua_rawgeti(l.s, C.LUA_REGISTRYINDEX, C.int(funi))
	if !isStackFunc(l.s, -1) {
		return false, nil
	}

	nargs := C.int(0)
	if args != nil {
		nargs = 1
		ext.Assert(C.lua_checkstack(l.s, nargs) != 0, "not enough free stack slots")
		parg := C.CString(string(args))
		defer C.free(unsafe.Pointer(parg))
		C.lua_pushlstring(l.s, parg, C.size_t(len(args)))
	}

	if C.lpcall(l.s, nargs, 1, 0, 0, nil) == C.LUA_OK {
	} else {
		return false, nil
	}
	retCnt := C.lua_gettop(l.s) - n
	if retCnt > 0 {
		retLength := C.lrawlen(l.s, 1)
		r := C.lua_tolstring(l.s, 1, &retLength)
		ret = C.GoBytes(unsafe.Pointer(r), C.int(retLength))
	}
	return true, ret
}

type LUserData interface {
	PushStack(l *C.lua_State)
}

//export TStt
type TStt struct {
	X   int
	Y   int
	Arr []int
}

var CLuaMetarArray = C.CString("luaMetaArray")

func (d *TStt) PushStack(l *C.lua_State) {
	rt := (*unsafe.Pointer)(C.lua_newuserdata(l, C.size_t(unsafe.Sizeof(d))))
	*rt = unsafe.Pointer(d)
	C.lua_getfield(l, C.LUA_REGISTRYINDEX, CLuaMetarArray)
	C.lua_setmetatable(l, -2)
}

//export TStt_new
func TStt_new(L unsafe.Pointer) C.int {
	l := (*C.lua_State)(L)
	//rt := unsafe.Pointer(C.lua_newuserdata(l, C.size_t(unsafe.Sizeof(TStt{}))))
	rt := (*unsafe.Pointer)(C.lua_newuserdata(l, C.size_t(unsafe.Sizeof(&TStt{}))))
	ptr := &TStt{}
	*rt = unsafe.Pointer(ptr)
	ptr.X = int(C.luaL_checkinteger(l, 1))
	ptr.Y = int(C.luaL_checkinteger(l, 2))
	ptr.Arr = make([]int, 3)
	ptr.Arr[0] = 5
	ptr.Arr[1] = ptr.X
	ptr.Arr[2] = 10
	C.lua_getfield(l, C.LUA_REGISTRYINDEX, CLuaMetarArray)
	C.lua_setmetatable(l, -2)
	return 1
}

//export TStt_setx
func TStt_setx(L unsafe.Pointer) C.int {
	l := (*C.lua_State)(L)
	a := (*unsafe.Pointer)(C.luaL_checkudata(l, 1, CLuaMetarArray))
	ptr := (*TStt)(*a)
	value := int(C.luaL_checkinteger(l, 2))

	ptr.X = value
	ptr.Arr[1] = value
	return 0
}

//export TStt_sety
func TStt_sety(L unsafe.Pointer) C.int {
	l := (*C.lua_State)(L)
	a := (*unsafe.Pointer)(C.luaL_checkudata(l, 1, CLuaMetarArray))
	ptr := (*TStt)(*a)
	value := int(C.luaL_checkinteger(l, 2))

	ptr.Y = value
	return 0
}

//export TStt_getx
func TStt_getx(L unsafe.Pointer) C.int {
	l := (*C.lua_State)(L)
	a := (*unsafe.Pointer)(C.luaL_checkudata(l, 1, CLuaMetarArray))
	ptr := (*TStt)(*a)

	C.lua_pushnumber(l, C.lua_Number(ptr.X))
	C.lua_pushnumber(l, C.lua_Number(ptr.Arr[1]))

	return 1
}

//export TStt_gety
func TStt_gety(L unsafe.Pointer) C.int {
	l := (*C.lua_State)(L)
	a := (*unsafe.Pointer)(C.luaL_checkudata(l, 1, CLuaMetarArray))
	ptr := (*TStt)(*a)

	//C.luaL_argcheck(l, NULL != a, 1, "'array' expected");

	C.lua_pushnumber(l, C.lua_Number(ptr.Y))

	return 1
}

//export TStt_gc
func TStt_gc(L unsafe.Pointer) C.int {
	//	l := (*C.lua_State)(L)
	//	a := (*unsafe.Pointer)(C.luaL_checkudata(l, 1, CLuaMetarArray))
	//	ptr := (*TStt)(*a)
	//*ptr = nil
	return 0
}
