package god

import (
	"bytes"
)

type Message interface{}

type Decoder interface {
	Decode(bytes.Buffer) *Message
}

type NodeDecoder struct {
}

type ClientDecoder struct {
}

type Encoder interface {
	Encode(*Message) bytes.Buffer
}

type NodeEncoder struct {
}

type ClientEncoder struct {
}

type Worker interface {
	ID() NID
	Stop()
}

type Runner interface {
	Run()
	Stop()
}

//func NewWorker(Runner) Worker

type receiveRunner struct {
}

//func NewReceiver(Decoder) Runner

type senderRunner struct {
}

//func NewSender(Encoder) Runner

type handleRunner struct{}

//func NewHandler(Handler) Runner

// func NewWorker(NewReceiver(NodeDecoder)) Runner

// Message in / out
type Handler interface {
	Handle(*Message)
}
