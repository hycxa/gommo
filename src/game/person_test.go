package game

import (
	"god"
	"testing"
)

func TestNewPerson(t *testing.T) {
	n1 := god.NewNode("n1", "tcp", "127.0.0.1:2008", god.NODE_GS_TYPE)
	p1 := interface{}(NewPerson(n1))
	_, ok := p1.(god.Handler)
	if !ok {
		t.Error("Not Handler")
	}
}
