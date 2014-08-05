package god

type MessageQueue interface {
	Push(source PID, target PID, m Message)   // noblocking
	Pop() (source PID, target PID, m Message) // blocking
}

type wholeMessage struct {
	source PID
	target PID
	Message
}

type messageList chan wholeMessage

type messageQueue struct {
	capacity int
	messageList
}

func NewMessageQueue(capacity int) MessageQueue {
	return &messageQueue{capacity, make(messageList, capacity)}
}

func (mq *messageQueue) Push(source PID, target PID, m Message) {
	mq.messageList <- wholeMessage{source, target, m}
}

func (mq *messageQueue) Pop() (source PID, target PID, m Message) {
	wm := <-mq.messageList
	return wm.source, wm.target, wm.Message
}
