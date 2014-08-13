package god

import (
	"ext"
	"net"
)

type acceptor struct {
	net.Listener
	NewAgent

	*stopper
}

func NewAcceptor(addr string, newAgent NewAgent) Stopper {
	listener, err := net.Listen("tcp", addr)
	ext.AssertE(err)
	a := &acceptor{Listener: listener, NewAgent: newAgent, stopper: NewStopper()}
	go ext.PCall(a.Run)
	return a
}

func (a *acceptor) Run() {
	defer a.Stopped()
	defer a.Listener.Close()

	for !a.StopRequested() {
		conn, err := a.Listener.Accept()
		if err == nil {
			a.NewAgent(conn)
		} else {
			ext.LogError(err)
			break
		}
	}
}
