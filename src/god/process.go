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
	Processor
	m Messenger
	h Handler
	PID
	observer PID
	mq       chan *proto.Message
	quit     chan int
}

func NewProcess(m Messenger, h Handler, observer PID) *process {
	o := new(process)
	o.PID.New()
	o.Handler = h
	o.mq = make(chan *proto.Message, CHAN_BUFF_NUM)
	o.quit = make(chan int)
	m.AddProcess(o)
	go o.run()
	return o
}

func (r *Process) run() {
	defer ext.UT(ext.T("Process::run"))
	for {
		select {
		case msg := <-r.mq:
			err := r.Handle(msg)
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

func (r *Process) pid() PID {
	return r.PID
}

func (r *Process) Notify(msg *proto.Message) error {
	r.mq <- msg
	return nil
}

func (r *Process) Call(msg *proto.Message) (error, *proto.Message) {
	r.mq <- msg
	return nil, nil
}
