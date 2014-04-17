package god

import (
	"proto"
	"sync"
)

type Messenger interface {
	Notify(PID, *proto.Message) error
	Call(PID, *proto.Message) (error, proto.Message)

	AddProcess(Processor)
	RemoveProcess(Processor)
	AllProcessInf(m *messenger) []PID

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
	ok, pro := m.processTab[PID]
	m.mutex.Unlock()

	if !ok {
		for i := 0; i < len(m.remoteTab); i++ {
			if m.remoteTab[i] != nil {
				if m.remoteTab[i].remoteNofity(msg) {
					return nil
				}
			}
		}
		return error.Error("not found process:%v", pid)
	}
	return pro.proNotify(msg)
}

func (m *messenger) Call(pid PID, msg *proto.Message) (error, *proto.Message) {
	m.mutex.Lock()
	ok, pro := m.processTab[PID]
	m.mutex.Unlock()
	if !ok {
		//TODO remoteCall
		return error.Error("not found process:%v", pid), nil
	}
	return pro.proCall(msg)
}

func (m *messenger) AddProcess(addObj Processor) {
	m.mutex.Lock()
	m.processTab[addObj.pid()] = addObj
	m.mutex.UnLock()

	for i := 0; i < len(m.remoteTab); i++ {
		if m.remoteTab[i] != nil {
			var msg Message
			msg.PacketID = proto.PROCESS_ADD_OR_REMOVE
			msg.Data = ProcessModify{UUID: addObj.pid(), IsAdd: true}
			m.remoteTab[i].remoteNofity(&msg)
		}
	}
}

func (m *messenger) RemoveProcess(delObj Processor) {
	m.mutex.Lock()
	delete(m.processTab, delObj.pid())
	m.mutex.UnLock()

	for i := 0; i < len(m.remoteTab); i++ {
		if m.remoteTab[i] != nil {
			var msg Message
			msg.PacketID = proto.PROCESS_ADD_OR_REMOVE
			msg.Data = ProcessModify{UUID: addObj.pid(), IsAdd: false}
			m.remoteTab[i].remoteNofity(&msg)
		}
	}
}

func (m *messenger) AddRemote(addObj *Remote) error {
	for i := 0; i < len(m.remoteTab); i++ {
		if m.remoteTab[i] == nil {
			m.remoteTab[i] = addObj
			return nil
		}
	}
}

func (m *messenger) RemoveRemote(delObj *Remote) {
	for i := 0; i < len(m.remoteTab); i++ {
		if m.remoteTab[i] == delObj {
			m.remoteTab[i] = nil
			return
		}
	}
}
