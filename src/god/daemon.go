package god

import (
	"ext"
	"net"
	"runtime"
	"sync"
)

type Daemon interface {
	net.Addr
}

type daemon struct {
	net.Listener
	net.Addr
}

var (
	nd   Daemon
	once sync.Once
)

func StartDaemon() Daemon {
	if nd != nil {
		return nd
	}

	once.Do(
		func() {
			d := &daemon{}
			ln, err := net.Listen("tcp", ":4369")

			ext.AssertE(err)

			d.Addr = ln.Addr()
			d.Listener = ln
			nd = d
		})
	for nd == nil {
		runtime.Gosched()
	}
	return nd
}
