package god

import (
	"bytes"
	"encoding/gob"
	"ext"
)

type nodeSender struct {
	Conn
	Encode
	Compress

	runner
	outgoing MessageQueue
}

func NewNodeSender(conn Conn, encode Encode, compress Compress) Runner {
	s := &nodeSender{Conn: conn, Encode: encode, Compress: compress}
	go s.Run()
	return s
}

func (s *nodeSender) Run() {
	for {
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
		var bufSize []byte
		BYTE_ORDER.PutUint32(bufSize, uint32(len(data)))
		s.Conn.Write(bufSize)
		s.Conn.Write(data)
	}
}
