package ext

import (
	"math"
	"runtime"
	"sync"
	"testing"
)

var (
	maxprocs = runtime.GOMAXPROCS(4)
)

type intBoolMap map[int]bool
type uint64BoolMap map[uint64]bool

func BenchmarkMapSet(b *testing.B) {
	m := make(intBoolMap)

	for i := 0; i < 100000; i++ {
		m[i] = true
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[i] = true
	}
}

func BenchmarkMapGet(b *testing.B) {
	m := make(intBoolMap)
	for i := 0; i < int(math.Max(float64(b.N), 100000)); i++ {
		m[i] = true
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Assert(m[i])
	}
}

func initParallelMap(m ParallelMap) {
	for i := 0; i < 100000; i++ {
		m.Set(i, i)
	}
}

func testParallelMap(t *testing.T, m ParallelMap) {
	initParallelMap(m)

	test := func(k, v interface{}) {
		AssertT(t, m.Set(k, v))
		ret, ok := m.Get(k).(uint64)
		AssertT(t, ok && v == ret)
		AssertT(t, m.Delete(k))
		AssertT(t, m.Get(k) == nil)
	}

	for i := -10000; i < 10000; i++ {
		v := RandomUint64()
		test(i, v)
	}

	for i := 0; i < 10000; i++ {
		k, v := string(RandomUint64()), RandomUint64()
		test(k, v)
	}
}

func TestChanMap(t *testing.T) {
	testParallelMap(t, NewChanMap())
}

func TestLockMap(t *testing.T) {
	testParallelMap(t, NewLockMap())
}

func TestLockFreeMap(t *testing.T) {
	testParallelMap(t, NewLockFreeMap())
}

func benchmarkParallelMap(b *testing.B, m ParallelMap) {
	var wg sync.WaitGroup
	wg.Add(4)
	initParallelMap(m)

	b.ResetTimer()
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			m.Set(string(RandomUint64()), true)
		}
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

	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			m.Len()
		}

	}()

	wg.Wait()
}

func benchmarkMapSet(b *testing.B, m ParallelMap) {
	initParallelMap(m)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			AssertB(b, m.Set(RandomUint64(), true))
		}
	})
}

func benchmarkMapGet(b *testing.B, m ParallelMap) {
	initParallelMap(m)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Get(RandomUint64())
		}
	})
}

func benchmarkMapDelete(b *testing.B, m ParallelMap) {
	initParallelMap(m)

	b.ResetTimer()
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
	benchmarkMapSet(b, NewChanMap())
}

func BenchmarkChanMapGet(b *testing.B) {
	benchmarkMapGet(b, NewChanMap())
}

func BenchmarkChanMapDelete(b *testing.B) {
	benchmarkMapDelete(b, NewChanMap())
}

func BenchmarkLockMap(b *testing.B) {
	benchmarkParallelMap(b, NewLockMap())
}

func BenchmarkLockMapSet(b *testing.B) {
	benchmarkMapSet(b, NewLockMap())
}

func BenchmarkLockMapGet(b *testing.B) {
	benchmarkMapGet(b, NewLockMap())
}

func BenchmarkLockMapDelete(b *testing.B) {
	benchmarkMapDelete(b, NewLockMap())
}
