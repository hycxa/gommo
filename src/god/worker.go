package god

import (
	"errors"
	"ext"
	"lua"
	"proto"
)

type Worker struct {
	*lua.L
	processFunRef C.int
	mq            chan *proto.Message
}

func NewWorker() *Worker {
	w := new(Worker)
	w.L = lua.NewLua()
	//TODO load lua file
	w.processFunRef = w.L.GetRef("Process")
	w.mq = make(chan proto.Message)
	go w.run()
	return w
}

func (w *Worker) run() {
	defer ext.UT(ext.T("worker::run"))
	defer w.L.Close()
	for {
		select {
		case m := <-w.mq:
			err := w.postMsg(m)
			if err != nil {
				ext.Errorf(err.Error())
			}
		}
	}
}

func (w *Worker) postMsg(data *proto.Message) error {
	ok, _ := w.L.CallRef(w.processFunRef, data.Data)
	if !ok {
		err := errors.New("Handle packetID:%v error", data.PackID)
		return err
	} else {
		return nil
	}
}

func (w *Worker) notify(data *proto.Message) {
	w.mq <- data
}

func (w *Worker) call() {

}
