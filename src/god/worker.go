package god

import (
	"ext"
)

var (
	workers = make(map[PID]Worker)
)

type worker struct {
	PID
	Handler
	Stopper
}

func NewWorker(id PID, h Handler) Worker {
	w := &worker{PID: id, Handler: h, Stopper: NewStopper()}
	go ext.PCall(w.run)
	return w
}

func (w *worker) Cast(source PID, msg Message) {
	w.Push(source, w.PID, msg)
}

func (w *worker) run() {
	defer w.Stopped()
	for !w.StopRequested() {
		source, target, msg := w.Pop()
		ext.Assert(source != target)
		// send to others
		if source == w.PID {
			w.BeforeSend(target, msg)
			w.Push(source, target, msg)
		} else {
			w.Handle(source, msg)
		}
	}
}
