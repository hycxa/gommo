package god

import (
	"errors"
	"lua"
)

type Worker struct {
	*lua.L
	processFunRef C.int
}

func NewWorker() Handler {
	w := new(Worker)
	w.L = lua.NewLua()
	//TODO load lua file
	w.processFunRef = w.L.GetRef("Process")
	return w
}

func (w *Worker) Handle(packID proto.PacketID, data *proto.Message) (error, []byte) {
	ok, rs := w.L.CallRef(w.processFunRef, data.Data)
	if !ok {
		err := errors.New("Handle packetID:%v error", packID)
		return err, nil
	} else {
		return nil, rs
	}
}

func (w *Worker) Close(){
	//TODO post luaState close operate
	w.L.Close
}

