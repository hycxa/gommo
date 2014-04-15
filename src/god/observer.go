package god

import (
	"proto"
)

type Observer interface {
	Add(interface{}) proto.UUID
	Remove(proto.UUID)
}
