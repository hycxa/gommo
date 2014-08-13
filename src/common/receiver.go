package common

import (
	"god"
	"net"
)

type receiver struct {
	god.Stopper
	handlerID god.PID
}

func NewReceiver(conn net.Conn, handlerID god.PID, decode god.Decode, decompress god.Decompress, decrypt god.Decrypt) god.Stopper {
	return &receiver{handlerID: handlerID, Stopper: god.NewStopper()}
}

func (r *receiver) Run() {
	defer r.Stopped()
}
