package proto

// name 全部大写
//结构体自动首字母大写

type LuaTransferData struct {
	bindata []byte
} /*
name:LUA_TRANSFER_DATA
desc: transfer lua data
*/

type Teq struct {
	X int
	Y int
} /*
name:XX1
desc:XXXXX

name:XX2
desc:XXXXXXX
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

type OtherUse struct {
	D int
	Q string
}

type Teq2 struct {
} /*
name:XX3
desc:XXXXX
*/

type OtherUse2 struct {
	D int
	Q string
}

type Teq3 struct {
} /*
name:XX4
desc:XXXXX
*/

type OtherUse3 struct {
	D int
	Q string
}
