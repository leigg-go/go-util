package _datastruct

import (
	"fmt"
	"sync"
)

/*
Two way linked list, suitable for frequent insertion / deletion scenarios.
*/

type node struct {
	value interface{}
	prev  *node
	next  *node
}

type LinkList struct {
	m     sync.RWMutex
	len   int
	first *node
	last  *node
}

func NewLinkList() *LinkList {
	return new(LinkList)
}

// Find is O(n), operation for elem is O(1).
func (l *LinkList) AddElemInBack(index int, v interface{}) {
	l.m.Lock()
	defer func() {
		l.len++
		l.m.Unlock()
	}()

	curr := l.first
	for i := 0; curr != nil; i++ {
		if i == index {
			newNode := &node{value: v, prev: curr, next: curr.next}
			curr.next = newNode
			return
		}
		curr = curr.next
	}
}

// Find is O(n), operation for elem is O(1).
func (l *LinkList) RemoveElem(v interface{}) error {
	l.m.Lock()
	defer func() {
		if l.first == nil {
			l.last = nil
		}
		l.len--
		l.m.Unlock()
	}()

	if l.len == 0 {
		return fmt.Errorf("elem [%v] is not exist", v)
	}

	curr := l.first
	for curr != nil {
		if curr.value == v {
			if curr.prev == nil {
				l.first = curr.next
			} else {
				curr.prev.next = curr.next
			}
			break
		}
		curr = curr.next
	}
	return nil
}

// GetElem only has find, O(n).
func (l *LinkList) GetElem(index int) interface{} {
	l.m.RLock()
	defer l.m.RUnlock()

	curr := l.first
	for i := 0; curr != nil; i++ {
		if i == index {
			return curr.value
		}
		curr = curr.next
	}
	return nil
}

func (l *LinkList) Clear() {
	l.m.Lock()
	defer l.m.Unlock()

	for l.first != nil {
		l.first = l.first.next
	}
	l.last = nil
	// The memory used by LinkList will be freed in next GC of runtime, not now.
	// but that is reusable, don't worry.
}
