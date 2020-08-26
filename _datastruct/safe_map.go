package _datastruct

import (
	"sync"
)

// ConcurrencyMap is a data container, which can be automatically de duplicated.
// ConcurrencyMap does not guarantee the order in which elements are stored.
type ConcurrencyMap struct {
	mu      sync.RWMutex
	h       map[interface{}]struct{}
	counter int
}

func NewConcurrencyMap() *ConcurrencyMap {
	return &ConcurrencyMap{
		mu: sync.RWMutex{},
		h:  make(map[interface{}]struct{}),
	}
}

// AddElem add elem to set, return true if not exists, return false otherwise.
func (s *ConcurrencyMap) AddElem(elem interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.h[elem]; ok {
		return false
	}
	s.h[elem] = struct{}{}
	s.counter++
	return true
}

// Contains return the given elem if contains in set.
func (s *ConcurrencyMap) Contains(elem interface{}) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.h[elem]; ok {
		return true
	}
	return false
}

// RemoveElem remove the elem from underlying map, return true if exists, return false otherwise.
func (s *ConcurrencyMap) RemoveElem(elem interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	exists := false
	if _, ok := s.h[elem]; ok {
		exists = true
		delete(s.h, elem)
		s.counter--
	}
	return exists
}

func (s *ConcurrencyMap) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.counter
}

// Range iterate underlying map.
func (s *ConcurrencyMap) Range(f func(int, interface{}) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	for k := range s.h {
		if !f(n, k) {
			return
		}
		n++
	}
}
