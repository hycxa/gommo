int lpcall(lua_State *L,
	int nargs,
	int nresults,
	int errfunc,
	int ctx,
	lua_CFunction k);
void lgetglobal(lua_State *L, const char *name);
lua_Number ltonumber(lua_State *L, int index);
size_t lrawlen (lua_State *L, int index);
#ifndef LUA_OK
#define LUA_OK 0
#endif
