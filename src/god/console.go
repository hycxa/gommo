package god

import (
	"runtime"
)

type Console struct {
}

func (c *Console) Run() {
	for {
		runtime.Gosched()
	}
}
