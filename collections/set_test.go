package collections

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestAdd(t *testing.T) {
	s := NewSet[int]()

	s.Add(1)
	s.Add(1)
	s.Add(2)
	s.Add(2)
	s.Add(3)

	assert.EqualValues(t, 3, s.Size())
}

func TestRemove(t *testing.T) {
	s := NewSet[int]()

	s.Add(1)
	s.Add(2)
	s.Remove(1)
	s.Remove(1)

	assert.EqualValues(t, 1, s.Size())
}

func TestContains(t *testing.T) {
	s := NewSet[int]()

	s.Add(1)
	s.Add(2)
	s.Remove(1)

	assert.EqualValues(t, s.Contains(1), false)
	assert.EqualValues(t, s.Contains(2), true)
}

func TestCollect(t *testing.T) {
	s := NewSet[int]()

	for i := 0; i < 1000; i++ {
		s.Add(i)
	}

	slice := s.Collect()
	sort.Ints(slice)
	for i := 0; i < 1000; i++ {
		assert.EqualValues(t, i, slice[i])
	}
}
