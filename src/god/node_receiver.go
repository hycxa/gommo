package god

import (
	"bytes"
	"encoding/gob"
	"ext"
	"io"
	"net"
)

type nodeReceiver struct {
	net.Conn
	Decode
	Decompress

	*runner
}

func NewNodeReceiver(conn net.Conn, decode Decode, decompress Decompress) Runner {
	return &nodeReceiver{Conn: conn, Decode: decode, Decompress: decompress, runner: NewRunner()}
}

func (r *nodeReceiver) Run() {
	defer r.Stopped()
	defer r.Conn.Close()

	ext.LogInfo("The connection [%s] accepted.", r.Conn.RemoteAddr().String())
	header := make([]byte, 4)

	for !r.StopRequested() {
		//r.Conn.SetReadDeadline(time.Now().Add(time.Minute))
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
