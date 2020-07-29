package _datastruct

import (
	"sync"
)

// HashSet is a data container, which can be automatically de duplicated.
// HashSet does not guarantee the order in which elements are stored.
type HashSet struct {
	mu sync.RWMutex
	h  map[interface{}]struct{}
}

func NewHashSet() *HashSet {
	return &HashSet{
		mu: sync.RWMutex{},
		h:  make(map[interface{}]struct{}),
	}
}

// AddElem add elem to set, return true if not exists, return false otherwise.
func (s *HashSet) AddElem(elem interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.h[elem]; ok {
		return false
	}
	s.h[elem] = struct{}{}
	return true
}

// Contains return the given elem if contains in set.
func (s *HashSet) Contains(elem interface{}) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.h[elem]; ok {
		return true
	}
	return false
}

// RemoveElem remove the elem from underlying map, return true if exists, return false otherwise.
func (s *HashSet) RemoveElem(elem interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	exists := false
	if _, ok := s.h[elem]; ok {
		exists = true
	}
	delete(s.h, elem)
	return exists
}

// Range iterate underlying map.
func (s *HashSet) Range(f func(int, interface{}) bool) {
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
