package god

import ()

type clientSender struct {
	runner
}

func NewClientSender(conn Conn, encode Encode, compress Compress, encrypt Encrypt) Runner {
	return nil
}

func (r *clientSender) Run() {
}
