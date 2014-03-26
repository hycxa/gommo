package god

import (
	"ext"
	"fmt"
	"proto"
)

var objects = make(map[proto.UUID]*Process)

type Process struct {
	Handler
	UUID proto.UUID
	mq   chan proto.Message
	quit chan int
}

func NewProcess(h Handler) *Process {
	o := new(Process)
	o.UUID.New()
	o.Handler = h
	o.mq = make(chan proto.Message)
	o.quit = make(chan int)

	objects[o.UUID] = o
	go o.run()
	return o
}

func Notify(source proto.UUID, target proto.UUID, packetID proto.PacketID, data interface{}) error {
	defer ext.UT(ext.T("NOTIFY"))
	t, ok := objects[target]
	if !ok {
		return fmt.Errorf("Target %v is not found!", target)
	}
	m := proto.Message{Sender: source, Data: data, PackID: packetID}
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
