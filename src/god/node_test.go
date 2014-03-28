package god

import (
	"ext"
	"testing"
	"time"
)

func TestNewNode(t *testing.T) {
	n1 := NewNode("n1", "tcp", "127.0.0.1:2001")
	n2 := NewNode("n2", "tcp", "127.0.0.1:2002")
	n3 := NewNode("n3", "tcp", "127.0.0.1:2003")
	n2.Dial("tcp", "127.0.0.1:2001")
	n3.DialNode(n2)
	n3.DialNode(n1)
	time.Sleep(100 * time.Millisecond)

	nodesOf1 := n1.Connected()

	on2, ok := nodesOf1["n2"]
	ext.AssertT(t, ok && n2.Name == on2.Name, "not found n2")

	on3, ok := nodesOf1["n3"]
	ext.AssertT(t, ok && n3.Name == on3.Name, "not found n2")
}
