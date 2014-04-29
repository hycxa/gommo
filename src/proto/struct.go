package proto

// name 全部大写
//结构体自动首字母大写
const (
	PACKAGE_SYSTEM = iota
	PACKAGE_USER
)

type CfgFlush struct {
	Db    string
	Url   string
	Modem string
} /*
name:CFG_FLUSH_REQ
packageScope:PACKAGE_SYSTEM
desc:flush config
*/

type CfgRsp struct {
	State bool
} /*
name:CFG_FLUSH_RSP
packageScope:PACKAGE_SYSTEM
desc: config respond
*/

type ProcessModify struct {
	UUID
	IsAdd bool
} /*
name:PROCESS_ADD_OR_REMOVE
packageScope:PACKAGE_SYSTEM
desc:process add or remove
*/
