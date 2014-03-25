package god

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	n1 := NewNode("n1", "tcp", "127.0.0.1:2001")
	n2 := NewNode("n2", "tcp", "127.0.0.1:2002")
	n3 := NewNode("n3", "tcp", "127.0.0.1:2003")
	n2.Dial("tcp", "127.0.0.1:2001")
	n3.DialNode(n2)

	nodesOf1 := n1.Connected()
	if on2, ok := nodesOf1["n2"]; !ok || n2.Name != on2.Name {
		t.Error("not found n2")
	}
	if on3, ok := nodesOf1["n3"]; !ok || n3.Name != on3.Name {
		t.Error("not found n3")
	}
}
