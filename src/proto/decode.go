package proto

import (
	"bytes"
	"encoding/gob"
)

func EncodeMsg(buff *bytes.Buffer, msg *Message) bool {
	enc := gob.NewEncoder(buff)
	err := enc.Encode(msg.Sender)
	if err != nil {
		checkErr(err)
		return false
	}
	err = enc.Encode(msg.PacketID)
	if err != nil {
		checkErr(err)
		return false
	}
	switch msg.PacketID {
	case LUA_TRANSFER_DATA:
		err = enc.Encode(msg.Data.(LuaTransferData))
	case CFG_FLUSH_REQ:
		err = enc.Encode(msg.Data.(CfgFlush))
	case CFG_FLUSH_RSP:
		err = enc.Encode(msg.Data.(CfgRsp))
	case PROCESS_ADD_OR_REMOVE:
		err = enc.Encode(msg.Data.(ProcessModify))
	default:
		return false
	}
	if err != nil {
		checkErr(err)
		return false
	}
	return true
}

func DecodeMsg(buff *bytes.Buffer) (bool, *Message) {
	msg := Message{}
	dec := gob.NewDecoder(buff)
	err := dec.Decode(&(msg.Sender))
	if err != nil {
		checkErr(err)
		return false, nil
	}
	err = dec.Decode(&(msg.PacketID))
	if err != nil {
		checkErr(err)
		return false, nil
	}
	switch msg.PacketID {
	case LUA_TRANSFER_DATA:
		var data LuaTransferData
		err = dec.Decode(&data)
		msg.Data = data
	case CFG_FLUSH_REQ:
		var data CfgFlush
		err = dec.Decode(&data)
		msg.Data = data
	case CFG_FLUSH_RSP:
		var data CfgRsp
		err = dec.Decode(&data)
		msg.Data = data
	case PROCESS_ADD_OR_REMOVE:
		var data ProcessModify
		err = dec.Decode(&data)
		msg.Data = data
	default:
		return false, nil
	}
	if err != nil {
		checkErr(err)
		return false, nil
	}
	return true, &msg
}
