package god

import (
	"errors"
	"ext"
	"proto"
	"testing"
)

func TestNewProessor(t *testing.T) {
	// o1 := NewProcess()
	// o1_b := GetNotifier(o1.ID)
	// if o1 != o1_b {
	// 	t.Error("GetNotifier")
	// }
}

type tv struct {
	i int
}

type process struct {
	*Process
	Name string
}

func new_process(name string) *process {
	p := new(process)
	p.Name = name
	p.Process = NewProcess(p)
	return p
}

func Handle(id proto.PacketID, m Marshaler) (retID proto.PacketID, ret Marshaler, err error) {
	switch m := m.(type) {
	case tv:
		v := tv(m)
		//v.i++
		return id, v, nil
	}
	return 0, nil, errors.New("wrong type")
}

func (p *process) Handle(id proto.PacketID, m Marshaler) (retID proto.PacketID, ret Marshaler, err error) {
	defer ext.UT(ext.T("process::Handle"))
	ext.Debugf("P[%s]%#v\n", p.Name, m)
	return Handle(id, m)
}

func TestProcess(t *testing.T) {
	if testing.Short() {
		return
	}

	p1 := new_process("p1")
	p2 := new_process("p2")

	for i := 0; i < 1; i++ {
		Notify(p1.UUID, p2.UUID, proto.PacketID(i), tv{i})
	}

}


