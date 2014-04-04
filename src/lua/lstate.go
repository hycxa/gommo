package lua

/*
//#cgo CFLAGS: -I/opt/local/include/
#cgo LDFLAGS: -llua
//-L/opt/local/lib
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

func NewL() *L {
	self := new(L)
	self.s = C.luaL_newstate()
	C.luaL_openlibs(self.s)
	return self
}

func (self *L) Close() {
	C.lua_close(self.s)
}

func (self *L) getRetValue(args C.int) *[]interface{} {
	ret := make([]interface{}, int(args))

	for i, index := C.int(1), 0; i <= args; i++ {
		t := C.lua_type(self.s, i)
		switch t {
		case C.LUA_TNIL:
			ret[index] = nil
		case C.LUA_TBOOLEAN:
			ret[index] = bool(C.lua_toboolean(self.s, i) != 0)
		case C.LUA_TNUMBER:
			ret[index] = int64(C.lua_tonumberx(self.s, i, nil))
		case C.LUA_TSTRING:
			ret[index] = C.GoString(C.lua_tolstring(self.s, i, nil))
		case C.LUA_TTABLE:
			ret[index] = nil
		case C.LUA_TFUNCTION:
			ret[index] = nil
		case C.LUA_TUSERDATA:
			ret[index] = nil
		case C.LUA_TTHREAD:
			ret[index] = nil
		case C.LUA_TLIGHTUSERDATA:
			ret[index] = nil
		default:
			ret[index] = nil
		}
		index++
	}
	return &ret
}

func (self *L) DoString(str string) (ret *[]interface{}) {

	oriCnt := C.lua_gettop(self.s)
	C.luaL_loadstring(self.s, C.CString(str))

	C.lua_pcallk(self.s, 0, C.LUA_MULTRET, 0, 0, nil)
	n := C.lua_gettop(self.s) - oriCnt
	ret = self.getRetValue(n)
	C.lua_settop(self.s, oriCnt)
	return
}

func (self *L) Call(f string, args ...interface{}) (ret *[]interface{}) {
	oriCnt := C.lua_gettop(self.s)
	C.lua_getglobal(self.s, C.CString(f))
	nargs := 0
	for _, v := range args {
		switch v.(type) {
		case int:
			C.lua_pushinteger(self.s, C.lua_Integer(v.(int)))
		case int64:
			C.lua_pushinteger(self.s, C.lua_Integer(v.(int64)))
		case string:
			C.lua_pushlstring(self.s, C.CString(v.(string)), C.size_t(len(v.(string))))
		case bool:
			if v.(bool) {
				C.lua_pushboolean(self.s, C.int(1))
			} else {
				C.lua_pushboolean(self.s, C.int(0))
			}
		default:
			C.lua_pushnil(self.s)
		}
		nargs++
	}

	C.lua_pcallk(self.s, C.int(nargs), C.LUA_MULTRET, 0, 0, nil)
	n := C.lua_gettop(self.s) - oriCnt
	ret = self.getRetValue(n)
	C.lua_settop(self.s, oriCnt)
	return
}
