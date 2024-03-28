package main

import (
	"slices"
	"sync"
)

type Set struct {
	vs map[int]struct{}
	mu sync.Mutex
}

func NewSet() *Set {
	return &Set{
		vs: make(map[int]struct{}),
	}
}

func (s *Set) Has(v int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.vs[v]
	return ok
}

func (s *Set) Add(v int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.vs[v]; ok {
		return
	}
	s.vs[v] = struct{}{}
}

func (s *Set) List() []int {
	s.mu.Lock()
	defer s.mu.Unlock()

	vs := make([]int, 0, len(s.vs))
	for k := range s.vs {
		vs = append(vs, k)
	}
	slices.Sort(vs)

	return vs
}
