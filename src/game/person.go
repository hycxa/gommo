package game

import (
	"god"
	"proto"
)

type Person struct {
	*god.Process
}

func NewPerson(node *god.Node) *Person {
	p := new(Person)
	p.Process = god.NewProcess(node, p)
	return p
}

func (p *Person) Handle(data *proto.Message) (err error) {
	return  nil
}
