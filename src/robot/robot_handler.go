package main

import (
	"god"
)

type clientHandler struct {
	*god.HandlerBase
}

func NewRobotHandler(id god.PID) god.Handler {
	self := &clientHandler{HandlerBase: god.NewHandlerBase(id)}
	return self
}

func (c *clientHandler) Handle(source god.PID, msg god.Message) {

}

func (c *clientHandler) BeforeSend(target god.PID, msg god.Message) {

}
