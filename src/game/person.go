package game

import (
	"god"
)

type Person struct {
	*god.Process
}

func NewPerson() *Person {
	p := new(Person)
	p.Process = god.NewProcess(p)
	return p
}

func (p *Person) Handle(pID god.PacketID, data god.Marshaler) (retID god.PacketID, ret god.Marshaler, err error) {
	return pID, data, nil
}
