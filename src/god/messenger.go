package god

import (
	"proto"
)

type Messenger interface {
	Notify(UUID, Message) (ok, error)
	Call(UUID, Message) (ok, Message)
}
