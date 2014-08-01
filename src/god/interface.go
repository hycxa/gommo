package god

import (
	"encoding/binary"
	"net"
)

type Header struct {
	Source PID
	Target PID
	Size   PID
}

type Message interface {
}

type MessageList chan Message

var (
	BYTE_ORDER = binary.LittleEndian
)

type Decode func([]byte) Message

func DefaultDecode([]byte) Message {
	return nil
}

type Encode func(Message) []byte

func DefaultEncode(Message) []byte {
	return nil
}

type Compress func([]byte) []byte

func DefaultCompress(in []byte) []byte {
	return in
}

type Decompress func([]byte) []byte

func DefaultDecompress(in []byte) []byte {
	return in
}

type Encrypt func([]byte) []byte
type Decrypt func([]byte) []byte

type Worker interface {
	PID() PID
	Cast(source PID, message Message)
	Stop()
}

type Runner interface {
	Run()
	Stop()
}

type Conn net.Conn
type NewAgent func(Conn)

type NodeSender interface {
	Cast(source PID, target PID, message Message)
}

type Messenger interface {
	Post(target PID, message Message)
}

// Message in / out
type Handler interface {
	Send(target PID, message Message)
	Handle(source PID, message Message)
}
