package god

func NewWorker(r Runner) Worker {
	return &worker{r}
}

type worker struct {
	Runner
}

func (w *worker) ID() NID {
	return 0
}
