package god

import (
	"ext"
	"net"
)

type Acceptor struct {
	net.Listener
	AgentCreator
}

func NewAcceptor(addr net.Addr, creator AgentCreator) Runner {
	listener, err := net.Listen("tcp", addr.String())
	ext.AssertE(err)
	return &Acceptor{listener, creator}
}

func (a *Acceptor) Run() {
	conn, err := a.Listener.Accept()
	ext.AssertE(err)

	a.Create(conn)
}

func (a *Acceptor) Stop() {

}
