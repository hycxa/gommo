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
	err := enc.Encode(msg.Sender.Sum(nil))
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
		err = enc.Encode(Teq2(msg.data))
	case XX3:
		err = enc.Encode(Teq3(msg.data))
	case XX4:
		err = enc.Encode(Teq3(msg.data))
	default:
		return false, nil
	}
	return true, buff
}
