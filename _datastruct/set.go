package _datastruct

import "sync"

type Set struct {
	mu sync.RWMutex
	mp map[interface{}]struct{}
}

func NewSet(size ...int) *Set {
	size = append(size, 0)
	return &Set{mp: make(map[interface{}]struct{}, size[0])}
}

func (s *Set) Add(item ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, i := range item {
		if _, ok := s.mp[i]; !ok {
			s.mp[i] = struct{}{}
		}
	}
}

func (s *Set) Remove(item ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, i := range item {
		delete(s.mp, i)
	}
}

func (s *Set) Contains(item interface{}) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.mp[item]
	return ok
}

func (s *Set) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.mp)
}

func (s *Set) ToStringList() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var ss []string
	for k := range s.mp {
		ss = append(ss, k.(string))
	}
	return ss
}

func (s *Set) ToIntList() []int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var ss []int
	for k := range s.mp {
		ss = append(ss, k.(int))
	}
	return ss
}

func (s *Set) ToUint32List() []uint32 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var ss []uint32
	for k := range s.mp {
		ss = append(ss, k.(uint32))
	}
	return ss
}

func (s *Set) ToInt64List() []int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var ss []int64
	for k := range s.mp {
		ss = append(ss, k.(int64))
	}
	return ss
}
