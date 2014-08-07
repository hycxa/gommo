package ext

import (
	"runtime"
	"sync"
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

func testParallelMap(t *testing.T, m ParallelMap) {
	AssertT(t, m.Set(10, 100))
	v, ok := m.Get(10).(int)
	AssertT(t, ok && v == 100)
	AssertT(t, m.Delete(10))
	AssertT(t, m.Get(10) == nil)
}

func TestChanMap(t *testing.T) {
	testParallelMap(t, NewChanMap())
}

func benchmarkParallelMap(b *testing.B, m ParallelMap) {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			m.Set(string(RandomUint64()), true)
		}
		//AssertB(b, m.Len() == b.N)
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			v := m.Get(string(RandomUint64()))
			if v == nil {
			}
		}

	}()

	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			m.Delete(string(RandomUint64()))
		}

	}()

	wg.Wait()
}

func benchmarkMapSet(b *testing.B, m ParallelMap) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			AssertB(b, m.Set(RandomUint64(), true))
		}
	})
}

func benchmarkMapGet(b *testing.B, m ParallelMap) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Get(RandomUint64())
		}
	})
}

func benchmarkMapDelete(b *testing.B, m ParallelMap) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			AssertB(b, m.Delete(RandomUint64()))
		}
	})
}

func BenchmarkChanMap(b *testing.B) {
	benchmarkParallelMap(b, NewChanMap())
}
func BenchmarkChanMapSet(b *testing.B) {
	m := NewChanMap()
	benchmarkMapSet(b, m)
}

func BenchmarkChanMapGet(b *testing.B) {
	m := NewChanMap()
	benchmarkMapGet(b, m)
}

func BenchmarkChanMapDelete(b *testing.B) {
	m := NewChanMap()
	benchmarkMapDelete(b, m)
}
