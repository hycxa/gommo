package god

import (
// "ext"
)

type handler struct {
	PID
	incoming MessageQueue
	outgoing MessageQueue
}

func NewHandler(self PID) *handler {
	return &handler{PID: self}
}

func (h *handler) Send(target PID, m *Message) {
	h.outgoing.Push(h.PID, target, m)
}

func (h *handler) Handle() {
	// source, target, message := h.incoming.Pop()
	// ext.Assert(target == h.PID)
	// h.Handle(source, message)
}
