package god

import (
	"errors"
	"ext"
	"time"
)

type stopper struct {
	stopped       chan bool
	stopRequested bool
}

func NewStopper() *stopper {
	return &stopper{stopped: make(chan bool, 1), stopRequested: false}
}

func (r *stopper) Stop() {
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

func (r *stopper) StopRequested() bool {
	return r.stopRequested
}

func (r *stopper) Stopped() {
	r.stopped <- true
}
