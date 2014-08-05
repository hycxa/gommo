package god

type clientReceiver struct {
	*runner
	handlerID PID
}

func NewClientReceiver(conn Conn, handlerID PID, decode Decode, decompress Decompress, decrypt Decrypt) Runner {
	return &clientReceiver{handlerID: handlerID, runner: NewRunner()}
}

func (r *clientReceiver) Run() {
	r.Stopped()
}
