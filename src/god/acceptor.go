package god

import (
	"ext"
	"net"
)

type acceptor struct {
	net.Listener
	NewAgent
}

func NewAcceptor(addr net.Addr, newAgent NewAgent) Runner {
	listener, err := net.Listen("tcp", addr.String())
	ext.AssertE(err)
	a := &acceptor{listener, newAgent}
	return a
}

func (a *acceptor) Run() {
	conn, err := a.Listener.Accept()
	ext.AssertE(err)
	a.NewAgent(conn)
}

func (a *acceptor) Stop() {

}
