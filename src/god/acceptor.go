package god

import (
	"ext"
	"net"
)

type acceptor struct {
	net.Listener
	NewAgent
	runner
}

func NewAcceptor(addr string, newAgent NewAgent) Runner {
	listener, err := net.Listen("tcp", addr)
	ext.AssertE(err)
	a := &acceptor{Listener: listener, NewAgent: newAgent}
	return a
}

func (a *acceptor) Run() {
	conn, err := a.Listener.Accept()
	ext.AssertE(err)
	a.NewAgent(conn)
}
