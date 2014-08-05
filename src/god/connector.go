package god

import (
	"ext"
	"net"
)

type connector struct {
	NewAgent

	workers WorkerMap
}

func NewConnector(newAgent NewAgent) Connector {
	return &connector{NewAgent: newAgent}
}

func (c *connector) Dial(address string) {
	conn, err := net.Dial("tcp", address)
	ext.AssertE(err)
	c.NewAgent(c.workers, conn)
}

func (c *connector) Stop() {
	for _, w := range c.workers {
		w.Stop()
	}
}
