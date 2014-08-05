package god

import (
	"errors"
	"ext"
	"time"
)

type runner struct {
	stopped       chan bool
	stopRequested bool
}

func NewRunner() *runner {
	return &runner{stopped: make(chan bool, 1), stopRequested: false}
}

func (r *runner) Stop() {
	r.stopRequested = true
	timeout := make(chan bool)
	go func() {
		time.Sleep(2 * time.Second)
		timeout <- true
	}()
	select {
	case <-r.stopped:
		return
	case <-timeout:
		ext.LogError(errors.New("Timeout when stopping!"))
	}
}

func (r *runner) StopRequested() bool {
	return r.stopRequested
}

func (r *runner) Stopped() {
	r.stopped <- true
}
