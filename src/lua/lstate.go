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

func (self *L) DoString(str string) (ret []interface{}) {

	oriCnt := C.lua_gettop(self.s)
	C.luaL_loadstring(self.s, C.CString(str))

	C.lua_pcallk(self.s, 0, C.LUA_MULTRET, 0, 0, nil)

	n := C.lua_gettop(self.s) - oriCnt
	ret = make([]interface{}, int(n))
	println(n)

	for i, index := C.int(1), 0; i <= n; i++ {
		t := C.lua_type(self.s, i)
		println("lua_type", t)
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
	C.lua_settop(self.s, oriCnt)

	return
}

func (self *L) Call(f string, v ...interface{}) (ret []interface{}) {
	ret = make([]interface{}, 0)
	return
}
