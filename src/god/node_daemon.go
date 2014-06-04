package god

import (
	"ext"
	"net"
)

type Daemon interface {
	net.Addr
}

type daemon struct {
	net.Listener
	net.Addr
}

var (
	nd  Daemon
	ndc = make(chan *daemon)
)

func StartNodeDaemon() Daemon {
	if nd != nil {
		return nd
	}

	go func() {
		d := &daemon{}
		ln, err := net.Listen("tcp", ":4369")

		if err != nil {
			ext.LogError(err)
			ndc <- nil
			return
		}

		d.Addr = ln.Addr()
		ndc <- d
	}()
	nd = <-ndc
	return nd
}
