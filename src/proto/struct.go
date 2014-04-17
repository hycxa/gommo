package proto

// name 全部大写
//结构体自动首字母大写

type LuaTransferData struct {
	bindata []byte
} /*
name:LUA_TRANSFER_DATA
desc: transfer lua data
*/

type CfgFlush struct {
	Db    string
	Url   string
	Modem string
} /*
name:CFG_FLUSH_REQ
desc:flush config
*/

type CfgRsp struct {
	State bool
} /*
name:CFG_FLUSH_RSP
desc: config respond
*/

type ProcessModify struct {
	UUID
	IsAdd bool
} /*
name:PROCESS_ADD_OR_REMOVE
desc:process add or remove
*/
