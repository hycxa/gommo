package god

var (
	workers = make(map[PID]Worker)
)

func NewWorker(r Runner) Worker {
	go r.Run()

	return &worker{Runner: r, pid: 0}
}

type worker struct {
	MessageQueue
	Runner
	pid PID
}

func (w *worker) PID() PID {
	return w.pid
}

func (w *worker) Cast(source PID, msg Message) {
	w.Push(source, w.pid, msg)
}
