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

func NewL() *L {
	l := new(L)
	l.s = C.luaL_newstate()
	C.luaL_openlibs(l.s)
	return l
}

func (self *L) Close() {
	C.lua_close(self.s)
}

func (self *L) DoString(str string) (ret []interface{}) {

	//int n = lua_gettop(m_impl->L);
	C.luaL_loadstring(self.s, C.CString(str))

	cn := C.lua_pcallk(self.s, 0, C.LUA_MULTRET, 0, 0, nil)

	println(cn)

	//int retCount = lua_gettop(m_impl->L) - n;
	n := int(cn)
	ret = make([]interface{}, n)

	for i := C.int(-1); ; i-- {
		t := C.lua_type(self.s, i)
		println("lua_type", t)
		switch t {
		case C.LUA_TNIL:

		case C.LUA_TBOOLEAN:
			ret[i] = bool(C.lua_toboolean(self.s, i) != 0)
		case C.LUA_TNUMBER:
			ret[i] = int64(C.lua_tonumberx(self.s, i, nil))
		case C.LUA_TSTRING:
			ret[i] = C.GoString(C.lua_tolstring(self.s, i, nil))
		}
	}

	return
}
