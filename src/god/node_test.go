package god

import (
	"ext"
	"testing"
	"time"
)

func sleepALittle() {
	time.Sleep(1000 * time.Millisecond)
}

func TestNewNode1(t *testing.T) {
	TestServer(t)
	NodeInit("n1", "tcp", "127.0.0.1", 8001)

	//for {
		//sleepALittle()
	//}
}

func TestConnect(t *testing.T) {
	n1 := NewNode("n1")
	n2 := NewNode("n2")

	n1.Connect("n2@127.0.0.1")

	ext.Assert(len(n2.Nodes()) == 1, "The node's amount is wrong!")
}
