package ext

import (
	"crypto/rand"
	"math"
	"math/big"
)

var (
	maxBigInt = big.NewInt(math.MaxInt64)
)

func RandomBigInt() *big.Int {
	n, err := rand.Int(rand.Reader, maxBigInt)
	AssertE(err)
	return n
}

func RandomUint64() uint64 {
	return RandomBigInt().Uint64()
}

func RandomInt64() int64 {
	return RandomBigInt().Int64()
}
