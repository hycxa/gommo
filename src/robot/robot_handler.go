package main

import (
	"god"
)

type clientHandler struct {
	god.Stopper
	god.Handler
}

func NewRobotHandler(sender god.PID) god.Runner {
	self := &clientHandler{Handler: god.NewHandler(0), Stopper: god.NewRunner()}
	return self
}

func (c *clientHandler) Handle(source god.PID, m *god.Message) {

}

func (c *clientHandler) BeforeOutgoing(target god.PID, m *god.Message) {

}

func (c *clientHandler) Run() {
	// for {
	// 	handler.Handle()
	// 	handler.outgoing
	// }
	defer c.Stopped()
}
