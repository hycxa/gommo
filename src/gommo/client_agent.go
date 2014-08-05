package main

import (
	"common"
	"god"
	"net"
)

func NewClientAgent(workers god.WorkerMap, conn net.Conn) {
	s := god.NewWorker(common.NewSender(conn, god.DefaultEncode, nil, nil))
	h := god.NewWorker(NewClientHandler(s.PID()))
	r := god.NewWorker(common.NewReceiver(conn, h.PID(), god.DefaultDecode, nil, nil))

	workers[s.PID()] = s
	workers[h.PID()] = h
	workers[r.PID()] = r
}
