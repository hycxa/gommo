package god

import (
	"net"
)

var (
	nodeManager map[PID]NodeSender
)

type Header struct {
	Source PID
	Target PID
	Size PID
}

type Message interface{
}

type Decoder interface {
	Decode([]byte) Message
}

type NodeDecoder struct {
}

type ClientDecoder struct {
}

type Encoder interface {
	Encode(Message) []byte
}

type NodeEncoder struct {
}

type ClientEncoder struct {
}

type Worker interface {
	PID() PID
	Cast(source PID, message Message)
	Stop()
}

func FindWorker(pid PID) Worker {
	
}

func FindNodeOfWorker(pid PID) NodeSender{
	
}

func Cast(source ID, target PID, message Message) {
	worker := FindWorker(target)
	if worker != nil {
		worker.Cast(source, message)
	}

	sender := FindNodeOfWorker(target)
	if sender != nil {
		sender.Cast(source, target, message)
	}
}

func (r *ClientReceiver) Run() {
	message := Decode()
	Cast(handlerID, handlerID, message)
}

type Runner interface {
	Run()
	Stop()
}

type AgentCreator interface {
	Create(net.Conn)
}

//func NewWorker(Runner) Worker

type Receiver struct {
}

//func NewReceiver(Decoder) Runner

type NodeSender interface {
	Cast(source, target, Message)
}

//func NewSender(Encoder) Runner

type handleRunner struct{}

//func NewHandler(Handler) Runner

// func NewWorker(NewReceiver(NodeDecoder)) Runner

// Message in / out
type Handler interface {
	Handle(source PID, Message)
}
