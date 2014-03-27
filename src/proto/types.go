package proto

import (
	"ext"
	"crypto/rand"
)

func checkErr(err error) {
	if err != nil {
		ext.Error(err)
	}
}

type PacketID int64

type UUID [16]byte

func (self *UUID) New() {
	var bt []byte = self[:]
	_, err := rand.Read(bt)
	if err != nil {
		checkErr(err)
		return
	}
}

type Message struct {
	Sender UUID
	Data   interface{}
	PackID PacketID
}
