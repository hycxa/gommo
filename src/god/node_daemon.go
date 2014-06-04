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

func StartNodeDaemon() Daemon {
	if nd != nil {
		return nd
	}

	go once.Do(
		func() {
			d := &daemon{}
			ln, err := net.Listen("tcp", ":4369")

			if err != nil {
				ext.LogError(err)
				return
			}

			d.Addr = ln.Addr()
			nd = d
		})
	for nd == nil {
		runtime.Gosched()
	}
	return nd
}
