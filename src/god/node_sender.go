package god

import (
	"bytes"
	"encoding/gob"
	"ext"
	"net"
	"time"
)

type nodeSender struct {
	net.Conn
	Encode
	Compress

	*runner
	outgoing MessageQueue
}

func NewNodeSender(conn net.Conn, encode Encode, compress Compress) Runner {
	s := &nodeSender{Conn: conn, Encode: encode, Compress: compress, runner: NewRunner(), outgoing: NewMessageQueue(32)}
	return s
}

func (s *nodeSender) Run() {
	defer s.Stopped()

	if RELEASE {
		ext.AssertE(s.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second)))
	}
	WriteBytes(s.Conn, GobEncode(MyInfo()))

	for !s.StopRequested() {
		source, target, msg := s.outgoing.Pop()
		data := s.Encode(msg)
		ext.Assert(data != nil)

		var buf bytes.Buffer

		enc := gob.NewEncoder(&buf)
		ext.AssertE(enc.Encode(source))
		ext.AssertE(enc.Encode(target))
		ext.AssertE(enc.Encode(len(data)))
		ext.AssertE(enc.Encode(data))

		data = buf.Bytes()
		if s.Compress != nil {
			data = s.Compress(buf.Bytes())
		}

		if RELEASE {
			ext.AssertE(s.Conn.SetWriteDeadline(time.Now().Add(time.Minute)))
		}
		WriteBytes(s.Conn, data)
	}
}
