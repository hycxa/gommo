package god

import (
	"net"
)

type ClentAgentCreator struct {
}

func (n *ClientAgentCreator) Create(conn net.Conn) {

	sender := NewSender(conn, &ClientEncoder{})
	NewWorker(sender)
	handler := 
	messenger := NewMessenger(NewClientHandler(sender))
	NewWorker()
	NewWorker(NewClientReceiver(conn, &ClientDecoder{}))
}

func (d *ClientDecoder) Decode([]byte) *Message {
	// sourceID, targetID = self()
	return nil
}

func (e *ClientEncoder) Encode(*Message) []byte {
	// sourceID, targetID = 0
	return nil
}

Worker -> Cast -> 

Sender 
	Cast(*Message)
	Run -> GetMessage -> ClientEncoder(targetID=0) -> Send(Message)

NodeSender
	Cast(*Message)

	>Run -> GetMessage -> NodeEncoder(do nothing) -> Send(Message)
	>
Messenger -> Run -> GetMessage -> Handle

ClientReceiver -> Run -> RecvSocket -> ClientDeccoder(sourceID=self()) -> GetWorker:Cast

NodeReceiver -> Run -> RecvSocket -> NodeDecoder(do nothing) -> GetWorker:Cast
