package proto

import (
	"crypto/rand"
	"ext"
)

func checkErr(err error) {
	if err != nil {
		ext.LogError(err)
	}
}

type PacketID int64

type UUID [16]byte

func (self *UUID) New() {
	var bt []byte = self[:]
	_, err := rand.Read(bt)
	if err != nil {
		ext.LogError(err)
		return
	}
}

type Message struct {
	Sender   UUID
	Reciever UUID
	Data     interface{}
	PackID   PacketID
}
