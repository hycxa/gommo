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
	err = enc.Encode(msg.PackID)
	if err != nil {
		checkErr(err)
		return false
	}
	switch msg.PackID {
	case LUA_TRANSFER_DATA:
		err = enc.Encode(msg.Data.(LuaTransferData))
	case XX1:
		err = enc.Encode(msg.Data.(Teq))
	case XX2:
		err = enc.Encode(msg.Data.(Teq))
	case CFG_FLUSH_REQ:
		err = enc.Encode(msg.Data.(CfgFlush))
	case CFG_FLUSH_RSP:
		err = enc.Encode(msg.Data.(CfgRsp))
	case XX3:
		err = enc.Encode(msg.Data.(Teq2))
	case XX4:
		err = enc.Encode(msg.Data.(Teq3))
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
	err = dec.Decode(&(msg.PackID))
	if err != nil {
		checkErr(err)
		return false, nil
	}
	switch msg.PackID {
	case LUA_TRANSFER_DATA:
		var data LuaTransferData
		err = dec.Decode(&data)
		msg.Data = data
	case XX1:
		var data Teq
		err = dec.Decode(&data)
		msg.Data = data
	case XX2:
		var data Teq
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
	case XX3:
		var data Teq2
		err = dec.Decode(&data)
		msg.Data = data
	case XX4:
		var data Teq3
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
