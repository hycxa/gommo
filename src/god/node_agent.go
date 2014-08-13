package god

import (
	"net"
)

func NewNodeAgent(conn net.Conn) {
	NewNodeSender(conn, DefaultEncode, nil)
	NewNodeReceiver(conn, DefaultDecode, nil)
}
