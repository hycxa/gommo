package god

import (
	"errors"
	"ext"
	"proto"
	"testing"
)

type tv struct {
	i int
}

type process struct {
	*Process
	Name string
}

func new_process(node *Node, name string) *process {
	p := new(process)
	p.Name = name
	p.Process = NewProcess(node, p)
	return p
}

func Handle(id proto.PacketID, m *proto.Message) (err error) {
	switch d := m.Data.(type) {
	case tv:
		v := tv(d)
		_ = v
		//v.i++
		return  nil
	}
	return errors.New("wrong type")
}

func (p *process) Handle(id proto.PacketID, m *proto.Message) (err error) {
	defer ext.UT(ext.T("process::Handle"))
	//ext.Debugf("P[%s]%#v\n", p.Name, m)
	return Handle(id, m)
}

func TestProcess(t *testing.T) {
	if testing.Short() {
		return
	}

	n1 := NewNode("n1", "tcp", "127.0.0.1:2008", NODE_GS_TYPE)
	p1 := new_process(n1, "p1")
	p2 := new_process(n1, "p2")

	for i := 0; i < 1; i++ {
		n1.Notify(p1.UUID, p2.UUID, proto.PacketID(i), tv{i})
	}

}
