package main

import (
	"god"
	"net"
)

type clientSender struct {
	god.Stopper
}

func NewRobotSender(conn net.Conn, encode god.Encode, compress god.Compress, encrypt god.Encrypt) god.Runner {
	return &clientSender{Stopper: god.NewRunner()}
}

func (s *clientSender) Run() {
	defer s.Stopped()
}
