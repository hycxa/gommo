package god

import (
	"ext"
)

type PID uint64
type NID uint64

func GeneratePID() PID {
	return PID(ext.RandomUint64())
}

func GenerateNID() NID {
	return NID(ext.RandomUint64())
}
