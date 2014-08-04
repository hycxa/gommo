package god

import (
	"ext"
)

type runner struct {
	stopRequest chan bool
}

func NewRunner() runner {
	return runner{make(chan bool, 2)}
}

func (r *runner) Stop() {
	defer ext.UT(ext.T())
	//r.stopRequest <- true
	//<-r.stopRequest
}
