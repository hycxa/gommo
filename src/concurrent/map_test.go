package concurrent

import (
	"ext"
	"testing"
)

type intBoolMap map[int]bool
type uint64BoolMap map[uint64]bool

func TestMapParallel(t *testing.T) {

}

func BenchmarkMapSet(b *testing.B) {
	m := make(intBoolMap)
	for i := 0; i < b.N; i++ {
		m[i] = true
	}
}

func BenchmarkMapGet(b *testing.B) {
	m := make(intBoolMap)
	for i := 0; i < b.N; i++ {
		m[i] = true
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ext.Assert(m[i])
	}
}

func BenchmarkMapParallelSet(b *testing.B) {
	m := make(uint64BoolMap)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m[ext.RandomUint64()] = true
		}
	})
}

func BenchmarkMapParallelGet(b *testing.B) {
	m := make(uint64BoolMap)

	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine has its own bytes.Buffer.
		for pb.Next() {
			_, ok := m[ext.RandomUint64()]
			ext.Assert(ok == false)
		}
	})
}
