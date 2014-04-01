package god

import (
	"ext"
	"testing"
	"time"
)

func sleepALittle() {
	time.Sleep(1000 * time.Millisecond)
}

func TestNewNode(t *testing.T) {
	TestServer(t)
	n1 := NewNode("n1", "tcp", "127.0.0.1:2001", NODE_GS_TYPE)
	n2 := NewNode("n2", "tcp", "127.0.0.1:2002", NODE_GATE_TYPE)
	n3 := NewNode("n3", "tcp", "127.0.0.1:2003", NODE_GS_TYPE)
	n4 := NewNode("n4", "tcp", "127.0.0.1:2004", NODE_GATE_TYPE)
	n1.ConnOtherSvr()
	n2.ConnOtherSvr()
	n3.ConnOtherSvr()
	n4.ConnOtherSvr()
	//n2.Dial("tcp", "127.0.0.1:2001")
	//n3.DialNode(n2)
	//n3.DialNode(n1)

	sleepALittle()

	nodesOf1 := n1.Connected()
	nodesOf2 := n2.Connected()
	nodesOf3 := n3.Connected()
	nodesOf4 := n4.Connected()

	for _, info := range nodesOf1 {
		ext.AssertT(t, info.Name == n2.Name || info.Name == n3.Name || info.Name == n4.Name, "has error conn", info.Name)
	}

	for _, info := range nodesOf2 {
		ext.AssertT(t, info.Name == n1.Name || info.Name == n3.Name, "has error conn", info.Name)
	}

	for _, info := range nodesOf3 {
		ext.AssertT(t, info.Name == n1.Name || info.Name == n2.Name || info.Name == n4.Name, "has error conn", info.Name)
	}

	for _, info := range nodesOf4 {
		ext.AssertT(t, info.Name == n1.Name || info.Name == n3.Name, "has error conn", info.Name)
	}
}
