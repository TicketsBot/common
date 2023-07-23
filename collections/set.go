package collections

import "encoding/json"

type Set[T comparable] struct {
	inner map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	set := Set[T]{
		inner: make(map[T]struct{}),
	}

	return &set
}

func (s *Set[T]) Add(value T) {
	s.inner[value] = struct{}{}
}

func (s *Set[T]) Remove(value T) {
	delete(s.inner, value)
}

func (s *Set[T]) Contains(value T) bool {
	_, ok := s.inner[value]
	return ok
}

func (s *Set[T]) Size() int {
	return len(s.inner)
}

func (s *Set[T]) Collect() []T {
	slice := make([]T, s.Size())
	i := 0
	for item := range s.inner {
		slice[i] = item
		i++
	}

	return slice
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Collect())
}

func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var slice []T
	if err := json.Unmarshal(data, &slice); err != nil {
		return err
	}

	for _, item := range slice {
		s.Add(item)
	}

	return nil
}
