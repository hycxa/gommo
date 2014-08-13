package god

import (
	"ext"
)

type HandlerBase struct {
	PID
	MessageQueue
}

func NewHandlerBase(id PID) *HandlerBase {
	return &HandlerBase{PID: id, MessageQueue: NewMessageQueue(32)}
}

func (h *HandlerBase) SendOut(target PID, msg Message) {
	ext.Assert(h.PID != target)
	h.Push(h.PID, target, msg)
}
