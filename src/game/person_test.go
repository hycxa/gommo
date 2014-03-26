package game

import (
	"god"
	"testing"
)

func TestNewPerson(t *testing.T) {
	p1 := interface{}(NewPerson())
	_, ok := p1.(god.Handler)
	if !ok {
		t.Error("Not Handler")
	}
}
