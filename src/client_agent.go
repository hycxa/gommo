package god

import ()

func NewClientAgent(workers WorkerMap, conn Conn) {
	s := NewWorker(NewClientSender(conn, DefaultEncode, nil, nil))
	h := NewWorker(NewClientHandler(s.PID()))
	r := NewWorker(NewClientReceiver(conn, h.PID(), DefaultDecode, nil, nil))

	workers[s.PID()] = s
	workers[h.PID()] = h
	workers[r.PID()] = r
}
