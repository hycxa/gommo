package god

import ()

type ClentNewAgent struct {
}

func NewClientAgent(conn Conn) {
	sender := NewWorker(NewClientSender(conn, DefaultEncode, nil, nil))
	handler := NewWorker(NewClientHandler(sender.PID()))
	NewWorker(NewClientReceiver(conn, handler.PID(), DefaultDecode, nil, nil))
}
