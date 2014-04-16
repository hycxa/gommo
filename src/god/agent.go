package god

import (
	"bytes"
	"ext"
	"io"
	"net"
)

type Agent struct {
	*process
	Messenger
	nFun     NotifyFun
	conn     net.Conn
	nodeAddr string
}

type agentHandler struct {
}

func NewAgent(m Messenger, nFun NotifyFun, conn net.Conn, nodeAddr string) Processor {
	a := new(Agent)
	a.process = NewProcess(m, new(agentHandler), nil)
	a.Messenger = m
	a.nFun = nFun
	a.conn = conn
	a.nodeAddr = nodeAddr
	go a.run()
	return a
}

func (a *Agent) notify(proto.Message) (ok, error) {
	a.notify(msg)
	return a.process.notify(msg)
}

func (a *Agent) call(proto.Message) (ok, proto.Message) {
	a.notify(msg)
	return a.process.call(msg)
}

func (a *Agent) run() {
	defer a.conn.Close()
	header := make([]byte, 2)

	for {
		//conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		_, err := io.ReadFull(a.conn, header)
		ext.AssertE(err)

		data := make([]byte, BYTE_ORDER.Uint16(header))
		//conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		_, err = io.ReadFull(a.conn, data)
		ext.AssertE(err)

		b := bytes.NewBuffer(data)
		ok, msg := proto.DecodeMsg(b)
		if ok {
			a.notify(msg)
		}
	}
}

func (a *agentHandler) Handle(packID proto.PacketID, data *proto.Message) error {
	return nil
}
