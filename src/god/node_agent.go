package god

import ()

func NewNodeAgent(workers WorkerMap, conn Conn) {
	r := NewWorker(NewNodeReceiver(conn, DefaultDecode, nil))
	s := NewWorker(NewNodeSender(conn, DefaultEncode, nil))
	workers[r.PID()] = r
	workers[s.PID()] = s
}
