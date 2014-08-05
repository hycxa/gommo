package god

import ()

type clientSender struct {
	*runner
}

func NewClientSender(conn Conn, encode Encode, compress Compress, encrypt Encrypt) Runner {
	return &clientSender{runner: NewRunner()}
}

func (s *clientSender) Run() {
	s.Stopped()
}
