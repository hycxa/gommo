package god

import (
	"bytes"
	"encoding/gob"
	"ext"
	"net"
	"time"
)

type nodeReceiver struct {
	net.Conn
	Decode
	Decompress

	*stopper
	nodeInfo
}

func NewNodeReceiver(conn net.Conn, decode Decode, decompress Decompress) Stopper {
	r := &nodeReceiver{Conn: conn, Decode: decode, Decompress: decompress, stopper: NewStopper()}
	go ext.PCall(r.Run)
	return r
}

func (r *nodeReceiver) Run() {
	defer r.Stopped()
	defer r.Conn.Close()

	if RELEASE {
		ext.AssertE(r.Conn.SetReadDeadline(time.Now().Add(10 * time.Second)))
	}
	infoData := ReadBytes(r.Conn)

	var info nodeInfo
	GobDecode(infoData, &info)

	ext.Assert(info.Cookie == MyInfo().Cookie)
	ext.Assert(info.ID != MyInfo().ID)

	r.nodeInfo = info
	AddNode(r.ID, r.nodeInfo)
	ext.LogDebug("ESTABLISHED\tLOCAL\t%s\tREMOTE\t%s\tINFO\t%v", r.LocalAddr().String(), r.RemoteAddr().String(), r.nodeInfo)

	for !r.StopRequested() {
		if RELEASE {
			ext.AssertE(r.Conn.SetReadDeadline(time.Now().Add(time.Minute)))
		}
		data := ReadBytes(r.Conn)
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
