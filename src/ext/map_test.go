package ext

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

// ChanMap size is 0, LockMap use RWMutex
// BenchmarkRandomUint64-4	 1000000	      1428 ns/op
// BenchmarkRandomInt64-4	 1000000	      1476 ns/op
// BenchmarkMapSet-4	10000000	       158 ns/op
// BenchmarkMapGet-4	50000000	        55.1 ns/op
// BenchmarkChanMap-4	  500000	      2732 ns/op
// BenchmarkChanMapSet-4	 5000000	       665 ns/op
// BenchmarkChanMapGet-4	 5000000	       741 ns/op
// BenchmarkChanMapDelete-4	 5000000	       581 ns/op
// BenchmarkLockMap-4	  500000	      3022 ns/op
// BenchmarkLockMapSet-4	 5000000	       665 ns/op
// BenchmarkLockMapGet-4	20000000	        92.9 ns/op
// BenchmarkLockMapDelete-4	 5000000	       543 ns/op

var (
	maxprocs          = runtime.GOMAXPROCS(4)
	initDatas         [1000000]uint64
	constInitDataSize = len(initDatas)
	initMapSize       = 100000
)

func init() {
	for i := 0; i < constInitDataSize; i++ {
		initDatas[i] = RandomUint64()
	}
}

type uint64BoolMap map[uint64]bool

func BenchmarkMapSet(b *testing.B) {
	m := make(uint64BoolMap)

	for i := 0; i < initMapSize; i++ {
		m[RandomUint64()] = true
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[initDatas[i%constInitDataSize]] = true
	}
}

func BenchmarkMapGet(b *testing.B) {
	m := make(uint64BoolMap)
	for i := 0; i < initMapSize; i++ {
		m[RandomUint64()] = true
	}
	AssertB(b, len(m) == initMapSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if m[initDatas[i%constInitDataSize]] {

		}
	}
}

func initParallelMap(m ParallelMap) {
	for i := 0; i < initMapSize; i++ {
		m.Set(i, i)
	}
}

func testParallelMap(t *testing.T, m ParallelMap) {
	initParallelMap(m)

	test := func(k, v interface{}) {
		m.Set(k, v)
		ret, ok := m.Get(k).(uint64)
		AssertT(t, ok && v == ret)
		m.Delete(k)
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

func benchmarkParallelMap(b *testing.B, m ParallelMap) {
	var wg sync.WaitGroup
	wg.Add(4)
	initParallelMap(m)

	b.ResetTimer()
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			m.Set(initDatas[i%constInitDataSize], true)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			v := m.Get(initDatas[i%constInitDataSize])
			if v == nil {
			}
		}

	}()

	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			m.Delete(initDatas[i%constInitDataSize])
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
			m.Set(time.Now().Nanosecond()%constInitDataSize, true)
		}
	})
}

func benchmarkMapGet(b *testing.B, m ParallelMap) {
	initParallelMap(m)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Get(time.Now().Nanosecond() % constInitDataSize)
		}
	})
}

func benchmarkMapDelete(b *testing.B, m ParallelMap) {
	initParallelMap(m)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Delete(initDatas[time.Now().Nanosecond()%constInitDataSize])
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
