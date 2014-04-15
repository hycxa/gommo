package god

import (
	"ext"
	"proto"
)

type PID proto.UUID

type Process struct {
	m Messenger
	h Handler
	PID
	observer PID
	mq       chan proto.Message
	quit     chan int
}

func NewProcess(m Messenger, h Handler, observer PID) *Process {
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
