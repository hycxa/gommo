package main

import (
	"common"
	"god"
	"net"
)

func NewRobotAgent(conn net.Conn) {
	common.NewSender(conn, god.DefaultEncode, nil, nil)
	id := god.GeneratePID()
	h := god.NewWorker(id, NewRobotHandler(id))
	common.NewReceiver(conn, h.ID(), god.DefaultDecode, nil, nil)
}
