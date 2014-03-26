package proto

import (
	"bytes"
	"encoding/gob"
	"errors"
	"ext"
	"god"
)

func checkErr(err error) {
	if err != nil {
		ext.Error(err)
	}
}

func EncodeMsg(msg *Message) (bool, bytes.Buffer) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(msg.Sender)
	if err != nil {
		checkErr(err)
		return false, nil
	}
	err = enc.Encode(msg.PackID)
	if err != nil {
		checkErr(err)
		return false, nil
	}
	switch msg.PackID {
	case XX1:
		err = enc.Encode(Teq(msg.data))
	case XX2:
		err = enc.Encode(Teq(msg.data))
	case XX3:
		err = enc.Encode(Teq2(msg.data))
	case XX4:
		err = enc.Encode(Teq3(msg.data))
	default:
		return false, nil
	}
	return true, buff
}

func DecodeMsg(buff *bytes.Buffer) (bool, Message) {
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
	case XX1:
		err = dec.Decode(&Teq(msg.data))
	case XX2:
		err = dec.Decode(&Teq(msg.data))
	case XX3:
		err = dec.Decode(&Teq2(msg.data))
	case XX4:
		err = dec.Decode(&Teq3(msg.data))
	default:
		return false, nil
	}
	return true, msg
}
func CreatePacketByPackID(packID PacketID) (PacketID, interface{}) {
	switch packID {
	case XX1:
		return packID, Teq{}
	case XX2:
		return packID, Teq{}
	case XX3:
		return packID, Teq2{}
	case XX4:
		return packID, Teq3{}
	default:
		return nil, nil
	}
}
