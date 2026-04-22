package frozenset

import (
	"iter"
	"maps"

	"github.com/swonky/set/internal/base"
)

var _ base.SetLike[int] = FrozenSet[int]{}
var _ base.AsSetter[int] = FrozenSet[int]{}

type FrozenSet[T comparable] struct {
	values map[T]struct{}
}

func New[T comparable](cap ...int) FrozenSet[T] {
	return FrozenSet[T]{values: make(map[T]struct{}, base.GetCap(cap...))}
}

func FromSetLike[S base.SetLike[T], T comparable](s S) FrozenSet[T] {
	switch v := any(s).(type) {
	case base.Set[T]:
		return FrozenSet[T]{values: maps.Clone(v)}
	case FrozenSet[T]:
		return FrozenSet[T]{values: maps.Clone(v.values)}
	}
	fs := New[T](s.Len())
	s.Range(
		func(t T) bool {
			fs.values[t] = struct{}{}
			return true
		},
	)
	return fs
}

func FromSlice[S ~[]T, T comparable](s S) FrozenSet[T] {
	fs := New[T](len(s))
	for _, t := range s {
		fs.values[t] = struct{}{}
	}
	return fs
}

func Collect[T comparable](seq iter.Seq[T]) FrozenSet[T] {
	s := New[T]()
	for t := range seq {
		s.values[t] = struct{}{}
	}
	return s
}

func Freeze[T comparable](sl base.SetLike[T]) FrozenSet[T] {
	r := make(base.Set[T], sl.Len())
	for t := range sl.Range {
		r.Add(t)
	}
	return FrozenSet[T]{values: r}

}

func (fs FrozenSet[T]) Contains(elem T) bool {
	_, ok := fs.values[elem]
	return ok
}

func (fs FrozenSet[T]) Len() int {
	return len(fs.values)
}

// Range
func (fs FrozenSet[T]) Range(yield func(T) bool) {
	for k := range fs.values {
		if !yield(k) {
			return
		}
	}
}

func (fs FrozenSet[T]) AsSet() base.Set[T] {
	return maps.Clone(fs.values)
}
