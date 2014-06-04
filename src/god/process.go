package god

import (
	"ext"
	"proto"
)

type PID proto.UUID

type Processor interface {
	pid() PID
	Notify(*proto.Message) error
	Call(*proto.Message) (error, *proto.Message)
}

type process struct {
	m    Messenger
	h    Handler
	p    proto.UUID
	mq   chan *proto.Message
	quit chan int
}

func NewProcess(m Messenger, h Handler) *process {
	o := new(process)
	o.p.New()
	o.h = h
	o.mq = make(chan *proto.Message, CHAN_BUFF_NUM)
	o.quit = make(chan int)
	m.AddProcess(o)
	go o.run()
	return o
}

func (r *process) run() {
	defer ext.UT(ext.T("Process::run"))
	for {
		select {
		case msg := <-r.mq:
			err := r.h.Handle(msg)
			if err != nil {
				ext.Errorf(err.Error())
			}
		case <-r.quit:
			return
		}
	}
}

func (r *process) Close() {
	r.quit <- 0
}

func (r *process) pid() PID {
	return PID(r.p)
}

func (r *process) Notify(msg *proto.Message) error {
	r.mq <- msg
	return nil
}

func (r *process) Call(msg *proto.Message) (error, *proto.Message) {
	r.mq <- msg
	return nil, nil
}
