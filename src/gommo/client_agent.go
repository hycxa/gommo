package main

import (
	"common"
	"god"
	"net"
)

func NewClientAgent(conn net.Conn) {
	common.NewSender(conn, god.DefaultEncode, nil, nil)
	id := god.GeneratePID()
	h := god.NewWorker(id, NewClientHandler(id))
	common.NewReceiver(conn, h.ID(), god.DefaultDecode, nil, nil)
}
