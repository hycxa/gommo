package god

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
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

func DefaultDecode(b []byte) Message {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	var m Message
	dec.Decode(&m)
	return m
}

type Encode func(Message) []byte

func DefaultEncode(m Message) []byte {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	enc.Encode(m)
	return buf.Bytes()
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

type Stopper interface {
	Stop()
	StopRequested() bool
	Stopped()
}

type Connector interface {
	Dial(address string)
	Stop()
}

type Runner interface {
	Run()
	Stopper
}

type Worker interface {
	PID() PID
	Cast(source PID, message Message)
	Runner
}

type WorkerMap map[PID]Worker

type NewAgent func(WorkerMap, net.Conn)

type NodeSender interface {
	Cast(source PID, target PID, message Message)
}

type Messenger interface {
	Post(target PID, message Message)
}

// Message in / out
type Handler interface {
	//Send(target PID, message Message)
	//	Handle(source PID, message Message)
}
