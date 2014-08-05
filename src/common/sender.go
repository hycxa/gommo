package common

import (
	"god"
	"net"
)

type sender struct {
	god.Stopper
}

func NewSender(conn net.Conn, encode god.Encode, compress god.Compress, encrypt god.Encrypt) god.Runner {
	return &sender{Stopper: god.NewRunner()}
}

func (s *sender) Run() {
	defer s.Stopped()
}
