package game

import (
	"god"
	"proto"
)

type Person struct {
	*god.Process
}

func NewPerson() *Person {
	p := new(Person)
	p.Process = god.NewProcess(p)
	return p
}

func (p *Person) Handle(pID proto.PacketID, data god.Marshaler) (retID proto.PacketID, ret god.Marshaler, err error) {
	return pID, data, nil
}
