package god

import (
	"proto"
)

type Messenger interface {
	Notify(PID, *proto.Message) (ok, error)
	Call(PID, *proto.Message) (ok, proto.Message)
	Add(Processor)
	Remove(Processor)
	AllProcessInfo() []PID
}
