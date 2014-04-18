package god

import (
	"ext"
	"net"
)

const (
	REMOTE_NODE_TYPE = iota
	CLIENT_TYPE
)

type Acceptor struct {
	net.Listener
	mes      Messenger
	selfType int
	laddr    string
}

func NewAcceptor(mes Messenger, selfType int, network string, laddr string) *Acceptor {
	acc := new(Acceptor)
	acc.selfType = selfType
	acc.laddr = laddr
	acc.mes = mes
	acc.Listener, err = net.Listen(network, laddr)
	if err != nil {
		ext.Errorf(err.Error())
		return nil
	}
	go acc.accept()
}

func (a *Acceptor) accept() {
	for {
		conn, err := a.Accept()
		if err != nil {
			ext.LogError(err)
		} else {
			if a.selfType == REMOTE_NODE_TYPE {
				NewRemote(m, conn)
			} else if a.selfType == CLIENT_TYPE {
				NewAgent(m, GetOneWorker(), conn)
			}
		}
	}
}
