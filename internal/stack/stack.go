package stack

import (
	"errors"
	"sync"
)

func New[K comparable]() *Stack[K] {
	return &Stack[K]{}
}

type (
	Stack[K comparable] struct {
		mu     sync.RWMutex
		values []K
		err    error
	}
)

func (s *Stack[K]) Push(k K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values = append(s.values, k)
}

func (s *Stack[K]) Pop() (k K, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.values) <= 0 {
		return
	}
	i := len(s.values) - 1
	k = s.values[i]
	ok = true
	s.values = s.values[:i]
	return
}

func (s *Stack[K]) Peek() (k K, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.values) <= 0 {
		return
	}
	k = s.values[len(s.values)-1]
	ok = true
	return
}

func (s *Stack[K]) Catch(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.err = errors.Join(s.err, err)
}

func (s *Stack[K]) Err() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.err
}

func (s *Stack[K]) Has(k K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.values {
		if v == k {
			return true
		}
	}
	return false
}

func (s *Stack[K]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.values)
}

func (s *Stack[K]) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values = nil
	s.err = nil
}
