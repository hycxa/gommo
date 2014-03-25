package god

import (
	// "bytes"
	"crypto/sha1"
	// "encoding/gob"
	"ext"
	"fmt"
)

var objects = make(map[UUID]*Process)

type Process struct {
	UUID
	Handler
	mq   chan Message
	quit chan int
}

func NewProcess(h Handler) *Process {
	o := new(Process)
	o.UUID = UUID{sha1.New()}
	o.Handler = h
	o.mq = make(chan Message)
	o.quit = make(chan int)

	objects[o.UUID] = o
	go o.run()
	return o
}

func Notify(source UUID, target UUID, packetID PacketID, data interface{}) error {
	defer ext.UT(ext.T("NOTIFY"))
	t, ok := objects[target]
	if !ok {
		return fmt.Errorf("Target %v is not found!", target)
	}
	m := Message{Sender: source, Data: data, PackID:packetID}
	t.mq <- m
	return nil
}

func (r *Process) run() {
	defer ext.UT(ext.T("Process::run"))
	for {
		select {
		case m := <-r.mq:
			retID, ret, err := r.Handle(m.PackID, m.Data)
			if err != nil {
				ext.Errorf(err.Error())
			}

			err = Notify(r.UUID, m.Sender, retID, ret)
			if err != nil {
				ext.Errorf(err.Error())
			}
		case <-r.quit:
			return
		}
	}
}

func (r *Process) Close() {
	r.quit <- 0
}
