#include <stdio.h>
#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>
#include "_cgo_export.h"

int openlibs(lua_State* L) {
 	luaL_openlibs(L);
 	return 0;
}

int newarray (lua_State *l){
	return TStt_new(l);	
}
int setx(lua_State *l) {
	return TStt_setx(l);
}

int sety(lua_State *l) {
	return TStt_sety(l);
}

int getx(lua_State *l) {
	return TStt_getx(l);
}

int gety(lua_State *l) {
	return TStt_gety(l);
}

int gc(lua_State *l) {
	return TStt_gc(l);
}

static const struct luaL_Reg arraylib_f[] = {
	{"new", newarray},
	{NULL, NULL}
};

static const struct luaL_Reg arraylib_m[] = {
	{"setx", setx},
	{"sety", sety},
	{"getx", getx},
	{"gety", gety},
	{"__gc", gc},
	{NULL, NULL}
};

void install_func(lua_State *L)
{
	luaL_newmetatable(L, "luaMetaArray");
	lua_pushstring(L, "__index");
	lua_pushvalue(L, -2);
	lua_settable(L, -3);

	#if LUA_VERSION_NUM == 502
	luaL_setfuncs(L, arraylib_m, 0);
	
	lua_newtable(L);
	luaL_newlib(L, arraylib_f);
	lua_setglobal(L, "array");
	#elif LUA_VERSION_NUM == 501
	luaL_register(L, NULL, arraylib_m);
	luaL_register(L, "array", arraylib_f);
	#endif
}

