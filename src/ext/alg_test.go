package ext

import (
	"testing"
)

func BenchmarkRandomUint64(b *testing.B) {
	m := make(map[uint64]bool)
	for i := 0; i < b.N; i++ {
		m[RandomUint64()] = true
	}
	AssertB(b, b.N == len(m), "random number dup")
}

func BenchmarkRandomInt64(b *testing.B) {
	m := make(map[int64]bool)
	for i := 0; i < b.N; i++ {
		m[RandomInt64()] = true
	}
	AssertB(b, b.N == len(m), "random number dup")
}
