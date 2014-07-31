package god

func NewWorker(r Runner) Worker {
	go r.Run()
	return &worker{r, 0}
}

type worker struct {
	Runner
	NID
}

func (w *worker) ID() NID {
	return w.NID
}

func (w *worker) () {
	
}
