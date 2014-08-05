package main

import (
	"god"
	"net"
)

type clientReceiver struct {
	god.Stopper
	handlerID god.PID
}

func NewRobotReceiver(conn net.Conn, handlerID god.PID, decode god.Decode, decompress god.Decompress, decrypt god.Decrypt) god.Runner {
	return &clientReceiver{handlerID: handlerID, Stopper: god.NewRunner()}
}

func (r *clientReceiver) Run() {
	defer r.Stopped()
}
