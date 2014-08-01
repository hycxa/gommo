package god

import ()

type NodeNewAgent struct {
}

func NewNodeAgent(conn Conn) {
	NewWorker(NewNodeReceiver(conn, DefaultDecode, nil))
	NewWorker(NewNodeSender(conn, DefaultEncode, nil))
}
