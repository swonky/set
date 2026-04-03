package frozenset

import (
	"iter"
	"maps"

	"github.com/swonky/set/internal/base"
)

type FrozenSet[T comparable] struct {
	smap map[T]struct{}
}

func New[T comparable](cap ...int) FrozenSet[T] {
	return FrozenSet[T]{smap: make(map[T]struct{}, base.GetCap(cap...))}
}

func From[T comparable](elems ...T) FrozenSet[T] {
	s := New[T](len(elems))
	for _, t := range elems {
		s.smap[t] = struct{}{}
	}
	return s
}

func Collect[T comparable](seq iter.Seq[T]) FrozenSet[T] {
	s := New[T]()
	for t := range seq {
		s.smap[t] = struct{}{}
	}
	return s
}

func (fs FrozenSet[T]) Contains(elem T) bool {
	_, ok := fs.smap[elem]
	return ok
}

func (fs FrozenSet[T]) Len() int {
	return len(fs.smap)
}

// Range
func (fs FrozenSet[T]) Range(yield func(T) bool) {
	for k := range maps.Clone(fs.smap) {
		if !yield(k) {
			return
		}
	}
}
