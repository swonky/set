package set

import (
	"iter"
	"maps"

	"github.com/swonky/set/internal/base"
	"github.com/swonky/set/types"
)

// New returns a new [Set].
// cap optionally specifies the initial capacity. If multiple values are
// provided, only the first is used.
func New[T comparable](cap ...int) Set[T] {
	return make(Set[T], base.GetCap(cap...))
}

// NewSync returns a new [SyncSet] backed by a builtin [Set].
//
// It is equivalent to [Wrap]([New][T](cap...)).
func NewSync[T comparable](cap ...int) *SyncSet[T] {
	return &SyncSet[T]{values: New[T](cap...)}
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

func FromSetLike[T comparable](s types.SetLike[T]) Set[T] {
	if e, ok := s.(Set[T]); ok {
		return maps.Clone(e)
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
