package set

import (
	"iter"
	"maps"

	"github.com/swonky/set/internal/base"
)

// New returns a new Set with a .
// If no items are provided, it returns an empty set.
func New[T comparable](cap ...int) Set[T] {
	return make(Set[T], base.GetCap(cap...))
}

// Collect returns a new Set containing all elements produced by the iterator.
// If a positive capacity hint is provided, it is used to preallocate the set.
func Collect[T comparable](it iter.Seq[T], size ...int) Set[T] {
	s := make(Set[T], base.GetCap(size...))
	for v := range it {
		s[v] = struct{}{}
	}
	return s
}

func FromSetLike[T comparable](s SetLike[T]) Set[T] {
	switch e := s.(type) {
	case Set[T]:
		return maps.Clone(e)
	case AsSetter[T]:
		return e.AsSet()
	}
	r := make(map[T]struct{}, s.Len())
	for v := range s.Range {
		r[v] = struct{}{}
	}
	return r
}

func FromSlice[S ~[]T, T comparable](s S) Set[T] {
	r := make(Set[T], len(s))
	for _, v := range s {
		r[v] = struct{}{}
	}
	return r
}
