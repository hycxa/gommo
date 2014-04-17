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
}

type agentHandler struct {
}

func NewAgent(m Messenger, nFun NotifyFun, conn net.Conn) Processor {
	a := new(Agent)
	a.process = NewProcess(m, new(agentHandler), nil)
	a.Messenger = m
	a.nFun = nFun
	a.conn = conn
	go a.run()
	return a
}

func (a *Agent) proNotify(msg *proto.Message) error {
	a.nFun.notify(msg)
	return a.process.proNotify(msg)
}

func (a *Agent) proCall(msg *proto.Message) (error, *proto.Message) {
	a.nFun.call(msg)
	return a.process.proCall(msg)
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
			a.nFun.notify(msg)
		}
	}
}

func (a *agentHandler) Handle(packID proto.PacketID, data *proto.Message) error {
	return nil
}
