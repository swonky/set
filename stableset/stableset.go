package stableset

import (
	"iter"
	"slices"

	"github.com/swonky/set"
	"github.com/swonky/set/internal/base"
)

type StableSet[T comparable] struct {
	values []T
	index  map[T]int
}

const promoteThreshold = 16

func (ss *StableSet[T]) Len() int {
	return len(ss.values)
}

func (ss *StableSet[T]) Range(yield func(T) bool) {
	for _, v := range ss.values {
		if !yield(v) {
			return
		}
	}
}

func (ss *StableSet[T]) Contains(elem T) bool {
	if ss.index != nil {
		_, ok := ss.index[elem]
		return ok
	}
	return slices.Contains(ss.values, elem)
}

func (ss *StableSet[T]) Add(elem T) {
	if ss.index != nil {
		if _, ok := ss.index[elem]; ok {
			return
		}
		ss.index[elem] = len(ss.values)
		ss.values = append(ss.values, elem)
		return
	}

	if slices.Contains(ss.values, elem) {
		return
	}

	ss.values = append(ss.values, elem)

	if len(ss.values) >= promoteThreshold {
		ss.promote()
	}
}

func (ss *StableSet[T]) Delete(elem T) {
	if ss.index != nil {
		i, ok := ss.index[elem]
		if !ok {
			return
		}

		delete(ss.index, elem)
		ss.values = slices.Delete(ss.values, i, i+1)

		for j := i; j < len(ss.values); j++ {
			ss.index[ss.values[j]] = j
		}
		return
	}

	if i := slices.Index(ss.values, elem); i >= 0 {
		ss.values = slices.Delete(ss.values, i, i+1)
	}
}

func (ss *StableSet[T]) promote() {
	m := make(map[T]int, len(ss.values))
	for i, v := range ss.values {
		m[v] = i
	}
	ss.index = m
}

func New[T comparable](cap ...int) *StableSet[T] {

	var (
		n  = base.GetCap(cap...)
		ss = &StableSet[T]{values: make([]T, 0)}
	)

	if n >= promoteThreshold {
		ss.index = make(map[T]int, n)
	}

	return ss
}

func Collect[T comparable](s iter.Seq[T]) *StableSet[T] {
	ss := New[T]()
	for v := range s {
		ss.Add(v)
	}
	return ss
}

func FromSetLike[T comparable](s set.SetLike[T]) *StableSet[T] {
	ss := New[T](s.Len())
	s.Range(func(t T) bool {
		ss.Add(t)
		return true
	})
	return ss
}

func FromSlice[S ~[]T, T comparable](s S) *StableSet[T] {
	n := len(s)
	if n == 0 {
		return New[T]()
	}

	ss := New[T](n)

	if n < promoteThreshold {
		for _, v := range s {
			if !slices.Contains(ss.values, v) {
				ss.values = append(ss.values, v)
			}
		}
		return ss
	}

	ss.index = make(map[T]int, n)

	for _, v := range s {
		if _, ok := ss.index[v]; ok {
			continue
		}
		ss.index[v] = len(ss.values)
		ss.values = append(ss.values, v)
	}

	return ss
}
