package god

import (
	"net"
)

func NewNodeAgent(workers WorkerMap, conn net.Conn) {
	h := NewWorker(NewNodeHandler(conn))
	r := NewWorker(NewNodeReceiver(conn, DefaultDecode, nil))
	s := NewWorker(NewNodeSender(conn, DefaultEncode, nil))
	workers[r.PID()] = r
	workers[s.PID()] = s
}
