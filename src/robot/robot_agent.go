package main

import (
	"god"
	"net"
)

func NewRobotAgent(workers god.WorkerMap, conn net.Conn) {
	s := god.NewWorker(NewRobotSender(conn, god.DefaultEncode, nil, nil))
	h := god.NewWorker(NewRobotHandler(s.PID()))
	r := god.NewWorker(NewRobotReceiver(conn, h.PID(), god.DefaultDecode, nil, nil))

	workers[s.PID()] = s
	workers[h.PID()] = h
	workers[r.PID()] = r
}
