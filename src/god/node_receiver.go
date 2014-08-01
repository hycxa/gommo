package god

import (
	"bytes"
	"encoding/gob"
	"ext"
	"io"
	"time"
)

type nodeReceiver struct {
	runner
	Conn
	Decode
	Decompress
}

func NewNodeReceiver(Conn Conn, decode Decode, decompress Decompress) Runner {
	return &nodeReceiver{Conn: Conn, Decode: decode, Decompress: decompress}
}

func (r *nodeReceiver) Run() {
	defer r.Conn.Close()
	header := make([]byte, 4)

	for {
		r.Conn.SetReadDeadline(time.Now().Add(time.Minute))
		_, err := io.ReadFull(r.Conn, header)
		ext.AssertE(err)

		data := make([]byte, BYTE_ORDER.Uint32(header))
		_, err = io.ReadFull(r.Conn, data)
		ext.AssertE(err)

		if r.Decompress != nil {
			data = r.Decompress(data)
		}

		buf := bytes.NewBuffer(data)

		dec := gob.NewDecoder(buf)

		var source, target PID
		var msgSize int
		ext.AssertE(dec.Decode(source))
		ext.AssertE(dec.Decode(target))
		ext.AssertE(dec.Decode(msgSize))

		msg := r.Decode(buf.Bytes())
		ext.Assert(msg != nil)

		Cast(source, target, msg)
	}
}
