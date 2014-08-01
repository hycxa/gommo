package god

type clientHandler struct {
	runner
	*handler
}

func NewClientHandler(sender PID) Runner {
	self := &clientHandler{handler: NewHandler(0)}
	return self
}

func (c *clientHandler) Handle(source PID, m *Message) {

}

func (c *clientHandler) BeforeOutgoing(target PID, m *Message) {

}

func (c *clientHandler) Run() {
	// for {
	// 	handler.Handle()
	// 	handler.outgoing
	// }
}
