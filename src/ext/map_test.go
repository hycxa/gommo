package ext

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

// ChanMap size is 0, LockMap use Mutex
// BenchmarkRandomUint64-4	 1000000	      1530 ns/op
// BenchmarkRandomInt64-4	 1000000	      1394 ns/op
// BenchmarkMapSet-4	10000000	       152 ns/op
// BenchmarkMapGet-4	50000000	        53.6 ns/op
// BenchmarkChanMap-4	  500000	      3297 ns/op
// BenchmarkChanMapSet-4	 1000000	      1012 ns/op
// BenchmarkChanMapGet-4	 5000000	       732 ns/op
// BenchmarkChanMapDelete-4	 5000000	       794 ns/op
// BenchmarkLockMap-4	 1000000	      1436 ns/op
// BenchmarkLockMapSet-4	 5000000	       628 ns/op
// BenchmarkLockMapGet-4	 5000000	       413 ns/op
// BenchmarkLockMapDelete-4	 5000000	       515 ns/op
//
// LockMap use RWMutex, total down, read up 4
// BenchmarkLockMap-4	  500000	      3125 ns/op
// BenchmarkLockMapSet-4	 5000000	       670 ns/op
// BenchmarkLockMapGet-4	20000000	        92.9 ns/op
// BenchmarkLockMapDelete-4	 5000000	       553 ns/op

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
	testParallelMap(t, newChanMap())
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
			AssertB(b, m.Set(time.Now().Nanosecond()%constInitDataSize, true))
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
			AssertB(b, m.Delete(initDatas[time.Now().Nanosecond()%constInitDataSize]))
		}
	})
}

func newChanMap() ParallelMap {
	return NewChanMap(true)
}

func BenchmarkChanMap(b *testing.B) {
	benchmarkParallelMap(b, newChanMap())
}
func BenchmarkChanMapSet(b *testing.B) {
	benchmarkMapSet(b, newChanMap())
}

func BenchmarkChanMapGet(b *testing.B) {
	benchmarkMapGet(b, newChanMap())
}

func BenchmarkChanMapDelete(b *testing.B) {
	benchmarkMapDelete(b, newChanMap())
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
