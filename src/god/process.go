package god

import (
	"ext"
	"proto"
)

type Process struct {
	Handler
	proto.UUID
	mq   chan proto.Message
	quit chan int
}

func NewProcess(node *Node, h Handler) *Process {
	o := new(Process)
	o.UUID.New()
	o.Handler = h
	o.mq = make(chan proto.Message)
	o.quit = make(chan int)

	node.AddProcess(o)
	go o.run()
	return o
}

func (r *Process) run() {
	defer ext.UT(ext.T("Process::run"))
	for {
		select {
		case m := <-r.mq:
			err := r.Handle(m.PackID, &m)
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
