package gogame

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

func (p *Person) Handle(data god.Marshaler) (ret god.Marshaler, err error) {
	return data, nil
}
