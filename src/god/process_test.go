package god

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"errors"
	"ext"
	"fmt"
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

func Handle(id PacketID, m Marshaler) (retID PacketID, ret Marshaler, err error) {
	switch m := m.(type) {
	case tv:
		v := tv(m)
		//v.i++
		return id, v, nil
	}
	return 0, nil, errors.New("wrong type")
}

func (p *process) Handle(id PacketID, m Marshaler) (retID PacketID, ret Marshaler, err error) {
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
		Notify(p1.UUID, p2.UUID, PacketID(i), tv{i})
	}

}

type A struct {
	Name string
	Data int
}

type MA struct {
	Data   interface{}
	PackID PacketID
}

func TestEncode(t *testing.T) {
	if testing.Short() {
		return
	}
	uuid := UUID{sha1.New()}
	_ = uuid
	dataA := A{Name: "abc", Data: 7}

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	dec := gob.NewDecoder(&buff)

	err := enc.Encode(1)
	if err != nil {
		fmt.Println("err enc", err)
	}
	err = enc.Encode(dataA)
	if err != nil {
		fmt.Println("err enc", err)
	}

	//err:=enc.Encode(Message{Sender: uuid, Data: dataA, PackID:1})
	//err:=enc.Encode(MA{Data:dataA, PackID:1})

	//var msg Message
	//var msg MA

	var packID PacketID
	err = dec.Decode(&packID)
	if err != nil {
		fmt.Println("err dec")
	}

	if packID == 1 {
		var msg A
		err = dec.Decode(&msg)
		if err != nil {
			fmt.Println("err dec")
		}
		fmt.Println("result is:", msg)
	}

}
