package god

import (
	"net"
)

func NewReceiver(net.Conn, Decoder) Runner {
	return &Receiver{}
}

func (r *NodeReceiver) Run() {

}

func (r *NodeReceiver) Stop() {

}

func (r *NodeReceiver) Run() {
	source, target, message := Decode()
	Cast(source, target, message)
}
