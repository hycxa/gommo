package god

import (
	"bytes"
)

type Worker interface {
	ID() NID
	Cast(bytes.Buffer)
}

func NewWorker(Handler) (Worker, NID) {
	return &worker{}, 0
}

type worker struct {
}

func (w *worker) ID() NID {
	return 0
}

func (w *worker) Cast(bytes.Buffer) {

}
