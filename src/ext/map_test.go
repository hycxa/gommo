package ext

import (
	//"math/rand"
	"runtime"
	//"sync"
	"testing"
)

var (
	maxprocs = runtime.GOMAXPROCS(4)
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
		Assert(m[i])
	}
}

// func BenchmarkMapParallel(b *testing.B) {
// 	m := make(map[string]bool)
// 	var wg sync.WaitGroup
// 	wg.Add(3)

// 	go func() {
// 		defer wg.Done()
// 		for {
// 			m[string(rand.Int())] = true
// 		}
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		for {
// 			if m[string(rand.Int())] {

// 			}
// 		}

// 	}()

// 	go func() {
// 		defer wg.Done()
// 		for {
// 			delete(m, string(rand.Int()))
// 		}

// 	}()

// 	wg.Wait()
// }

// func BenchmarkMapParallelSet(b *testing.B) {
// 	m := make(uint64BoolMap)

// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			m[ext.RandomUint64()] = true
// 		}
// 	})
// }

// func BenchmarkMapParallelGet(b *testing.B) {
// 	m := make(uint64BoolMap)

// 	b.RunParallel(func(pb *testing.PB) {
// 		// Each goroutine has its own bytes.Buffer.
// 		for pb.Next() {
// 			_, ok := m[ext.RandomUint64()]
// 			ext.Assert(ok == false)
// 		}
// 	})
// }
