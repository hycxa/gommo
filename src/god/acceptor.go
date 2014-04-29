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
	var err error
	acc.Listener, err = net.Listen(network, laddr)
	if err != nil {
		ext.Errorf(err.Error())
		return nil
	}
	go acc.accept()
	return acc
}

func (acc *Acceptor) accept() {
	for {
		conn, err := acc.Accept()
		if err != nil {
			ext.LogError(err)
		} else {
			if acc.selfType == REMOTE_NODE_TYPE {
				NewRemote(acc.mes, conn)
			} else if acc.selfType == CLIENT_TYPE {
				NewAgent(acc.mes, GetOneWorker(), conn)
			}
		}
	}
}
