package proto

import (
	"bytes"
	"fmt"
	"testing"
)

func TestUse(t *testing.T) {
	if testing.Short() {
		return
	}
	packetID := XX1
	var packet Teq
	packet.X = 5
	packet.Y = 6

	var buff bytes.Buffer
	msg := Message{}
	msg.Sender.New()
	msg.PacketID = PacketID(packetID)
	msg.Data = packet
	fmt.Println(msg)

	ret := EncodeMsg(&buff, &msg)
	if ret == false {
		fmt.Println("err enc")
		return
	}

	ret, rMsg := DecodeMsg(&buff)
	if ret == false {
		fmt.Println("err dec")
		return
	}
	fmt.Println(*rMsg)
}
