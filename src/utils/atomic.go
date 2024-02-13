package utils

import "sync"

type LockedMapBool struct {
	lock sync.Mutex
	M    map[string]bool
}

func (m *LockedMapBool) Get(k string) (bool, bool) {
	m.lock.Lock()
	v, ok := m.M[k]
	m.lock.Unlock()
	return v, ok
}

func (m *LockedMapBool) Set(k string, v bool) {
	m.lock.Lock()
	m.M[k] = v
	m.lock.Unlock()
	return
}

var data = LockedMapBool{M: map[string]bool{}}
