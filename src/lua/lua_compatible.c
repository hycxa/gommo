#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>

int lpcall(lua_State *L,
	int nargs,
	int nresults,
	int errfunc,
	int ctx,
	lua_CFunction k)
{
	#if LUA_VERSION_NUM == 502
	return lua_pcallk(L, nargs, nresults, errfunc, ctx, k);
	#elif LUA_VERSION_NUM == 501
	return lua_pcall(L, nargs, nresults, errfunc);
	#endif
}

void lgetglobal(lua_State *L, const char *name)
{
	lua_getglobal(L, name);
}

lua_Number ltonumber(lua_State *L, int index)
{
	return lua_tonumber (L, index);
}

size_t lrawlen (lua_State *L, int index)
{
	#if LUA_VERSION_NUM == 502
	return lua_rawlen (L, index);
	#elif LUA_VERSION_NUM == 501
	return lua_objlen (L, index);
	#endif
}
