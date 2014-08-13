package god

import (
	"encoding/binary"
	"net"
)

var (
	DEFAULT_BYTE_ORDER = binary.LittleEndian
)

type Message interface {
}

type Decode func([]byte) Message
type Encode func(Message) []byte

type Compress func([]byte) []byte
type Decompress func([]byte) []byte

type Encrypt func([]byte) []byte
type Decrypt func([]byte) []byte

type Stopper interface {
	Stop()
	StopRequested() bool
	Stopped()
}

type Connector interface {
	Dial(address string)
}

type Worker interface {
	ID() PID
	Cast(source PID, msg Message)
	Stopper
}

type WorkerMap map[PID]Worker

type NewAgent func(net.Conn)

type NodeSender interface {
	Cast(source PID, target PID, msg Message)
}

type Messenger interface {
	Post(target PID, msg Message)
}

// Message in / out
type Handler interface {
	SendOut(target PID, msg Message) // Handler 内部用于发出消息时使用
	BeforeSend(target PID, msg Message)
	Handle(source PID, msg Message)
	MessageQueue
}
