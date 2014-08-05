package god

import (
	"ext"
	"net"
)

type acceptor struct {
	net.Listener
	NewAgent

	*runner
	workers WorkerMap
}

func NewAcceptor(addr string, newAgent NewAgent) Runner {
	listener, err := net.Listen("tcp", addr)
	ext.AssertE(err)
	a := &acceptor{Listener: listener, NewAgent: newAgent, runner: NewRunner(), workers: make(WorkerMap)}
	return a
}

func (a *acceptor) Run() {
	defer a.Stopped()
	for !a.StopRequested() {
		conn, err := a.Listener.Accept()
		if err == nil {
			a.NewAgent(a.workers, conn)
		} else {
			ext.LogError(err)
			break
		}
	}
}

func (a *acceptor) Stop() {
	a.Listener.Close()
	for _, w := range a.workers {
		w.Stop()
	}
	a.runner.Stop()
}
