package god

type clientReceiver struct {
	handlerID PID
}

func NewClientReceiver(conn Conn, handlerID PID, decode Decode, decompress Decompress, decrypt Decrypt) Runner {
	return nil
}

func (r *clientReceiver) Run() {
}

func (c *clientReceiver) Stop() {

}
