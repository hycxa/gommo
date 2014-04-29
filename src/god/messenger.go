package god

import (
	"ext"
	"proto"
	"sync"
)

type Messenger interface {
	Notify(PID, *proto.Message) error
	Call(PID, *proto.Message) (error, *proto.Message)

	AddProcess(Processor)
	RemoveProcess(Processor)
	AllProcessInfo() []PID

	AddRemote(*Remote) error
	RemoveRemote(*Remote)
}

type messenger struct {
	processTab map[PID]Processor
	mutex      sync.Mutex
	remoteTab  []*Remote
}

func NewMessenger() Messenger {
	m := new(messenger)
	m.remoteTab = make([]*Remote, 1000)
	return m
}

func (m *messenger) Notify(pid PID, msg *proto.Message) error {
	m.mutex.Lock()
	pro, ok := m.processTab[pid]
	m.mutex.Unlock()

	if !ok {
		for i := 0; i < len(m.remoteTab); i++ {
			if m.remoteTab[i] != nil {
				if m.remoteTab[i].remoteNofity(msg) {
					return nil
				}
			}
		}
		return ext.MyError{"not found process"}
	}
	return pro.Notify(msg)
}

func (m *messenger) Call(pid PID, msg *proto.Message) (error, *proto.Message) {
	m.mutex.Lock()
	pro, ok := m.processTab[pid]
	m.mutex.Unlock()
	if !ok {
		//TODO remoteCall
		return ext.MyError{"not found process"}, nil
	}
	return pro.Call(msg)
}

func (m *messenger) AddProcess(addObj Processor) {
	m.mutex.Lock()
	m.processTab[addObj.pid()] = addObj
	m.mutex.Unlock()

	for i := 0; i < len(m.remoteTab); i++ {
		if m.remoteTab[i] != nil {
			var msg proto.Message
			msg.PacketID = proto.PROCESS_ADD_OR_REMOVE
			msg.Data = proto.ProcessModify{UUID: proto.UUID(addObj.pid()), IsAdd: true}
			m.remoteTab[i].remoteNofity(&msg)
		}
	}
}

func (m *messenger) RemoveProcess(delObj Processor) {
	m.mutex.Lock()
	delete(m.processTab, delObj.pid())
	m.mutex.Unlock()

	for i := 0; i < len(m.remoteTab); i++ {
		if m.remoteTab[i] != nil {
			var msg proto.Message
			msg.PacketID = proto.PROCESS_ADD_OR_REMOVE
			msg.Data = proto.ProcessModify{UUID: proto.UUID(delObj.pid()), IsAdd: false}
			m.remoteTab[i].remoteNofity(&msg)
		}
	}
}

func (m *messenger) AllProcessInfo() []PID {
	ret := make([]PID, 1)
	m.mutex.Lock()
	for pid, _ := range m.processTab {
		ret = append(ret, pid)
	}
	m.mutex.Unlock()
	return ret
}

func (m *messenger) AddRemote(addObj *Remote) error {
	for i := 0; i < len(m.remoteTab); i++ {
		if m.remoteTab[i] == nil {
			m.remoteTab[i] = addObj
			return nil
		}
	}
	return ext.MyError{"Has no empty remotetab"}
}

func (m *messenger) RemoveRemote(delObj *Remote) {
	for i := 0; i < len(m.remoteTab); i++ {
		if m.remoteTab[i] == delObj {
			m.remoteTab[i] = nil
			return
		}
	}
}
