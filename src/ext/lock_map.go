package ext

import (
	"sync"
)

type lockMap struct {
	m     map[interface{}]interface{}
	mutex sync.Mutex
}

func NewLockMap() ParallelMap {
	l := lockMap{}
	l.m = make(map[interface{}]interface{})
	return &l
}

func (l *lockMap) Set(k, v interface{}) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.m[k] = v
	return true
}

func (l *lockMap) Get(k interface{}) interface{} {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	v, ok := l.m[k]
	if ok {
		return v
	} else {
		return nil
	}
}

func (l *lockMap) Delete(k interface{}) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	delete(l.m, k)
	return true
}

func (l *lockMap) Len() int {
	return len(l.m)
}
