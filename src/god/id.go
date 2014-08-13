package god

import (
	"ext"
)

type PID uint64

func (id PID) ID() PID {
	return id
}

type NID uint64

func (id NID) ID() NID {
	return id
}

func GeneratePID() PID {
	return PID(ext.RandomUint64())
}

func GenerateNID() NID {
	return NID(ext.RandomUint64())
}
