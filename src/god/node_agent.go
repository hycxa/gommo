package god

import (
	"net"
)

type NodeAgentCreator struct {
}

func (n *NodeAgentCreator) Create(conn net.Conn) {
	NewWorker(NewReceiver(conn, &NodeDecoder{}))
	NewWorker(NewSender(conn, &NodeEncoder{}))
}

func (d *NodeDecoder) Decode([]byte) *Message {
	return nil
}

func (e *NodeEncoder) Encode(*Message) []byte {
	return nil
}
