type NID uint64 // unique number identifier

type Node struct {
	acceptor 	Worker
	nodeManager Worker
	receiverManager Worker
	Stop()
}
// NewNode(type, listen) *Node, NID

type Worker interface {
	Run()
	Stop()
	ID() NID
	Type() string
	Call(binary)
	Cast(binary)
}

type worker struct {
	handler Handler
	Worker
}
// NewWorker(type, handler) Worker, NID

type Handler interface {
	HandleInfo(binary)
	HandleCall()
	HandleCast()
}
