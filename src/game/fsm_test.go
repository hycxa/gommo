package gogame

import ()

type world struct {
	State
	c chan uint64
}

func Run() uint64 {
	//var w world = world{State: NewState(), c: make(chan uint64)}
	//message_id := <-w.c
	return 0
}
