package god

import (
	"crypto/rand"
	"ext"
	"math"
	"math/big"
)

type PID uint64

func GeneratePID() PID {
	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	ext.AssertE(err)
	return PID(n.Uint64())
}
