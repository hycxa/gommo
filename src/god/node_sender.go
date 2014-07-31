package god

import (
	"net"
)

type nodeMessage struct {
	source PID
	target PID
	Message
}

type nodeSender struct {
	net.Conn
	messageList chan[nodeMessage]
}

func NewNodeSender(net.Conn, Encoder) Runner {
	n := &nodeSender{}
	n.messageList = {}
	nodeManager.Add(pid, nodeSender)
	return &nodeSender{}
}

func (s *nodeSender) Cast(source PID, target PID, message Message) {

}
func (r *nodeSender) Run() {
	m <- r.messageList
	var h Header
	h.Source = m.source
	h.Target = m.target
	data := Encode(m)
	h.Size = len(data)
	headerBinary = Encode(h)
	size := sizeof(m.source) + sizeof(m.target) + len(bin)
	r.Conn.Write(size)
	r.Conn.Write(m.source)
	r.Conn.Write(m.target)
	r.Conn.Write(m)
}

func (r *nodeSender) Stop() {

}
