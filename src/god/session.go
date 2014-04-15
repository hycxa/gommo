package god

import (
	"bytes"
	"ext"
	"io"
	"net"
)

type Session struct {
	Messenger
	wpid     PID
	conn     net.Conn
	nodeAddr string
}

func NewSession(m Messenger, wpid PID, conn net.Conn, nodeAddr string) Handler {
	s := new(Session)
	s.Messenger = m
	s.wpid = wpid
	s.conn = conn
	s.nodeAddr = nodeAddr
	go s.run()
	return s
}

func (s *Session) run() {
	header := make([]byte, 2)

	for {
		//conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		_, err := io.ReadFull(s.conn, header)
		ext.AssertE(err)

		data := make([]byte, BYTE_ORDER.Uint16(header))
		//conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		_, err = io.ReadFull(s.conn, data)
		ext.AssertE(err)

		b := bytes.NewBuffer(data)
		ok, msg := proto.DecodeMsg(b)
		if ok {
			s.Handle(msg.PackID, msg)
		}
	}
}

func (s *Session) Close() {
	s.conn.Close()
}

func (s *Session) Handle(packID proto.PacketID, data *proto.Message) (error, []byte) {
	if packID == proto.LUA_TRANSFER_DATA {
		s.Notify(s.wpid, msg)
	} else {
	//do other self thing
	}
}
