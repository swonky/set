package frozenset

import (
	"iter"
	"maps"

	"github.com/swonky/set"
	"github.com/swonky/set/internal/base"
	"github.com/swonky/set/types"
)

func New[T comparable](cap ...int) FrozenSet[T] {
	return FrozenSet[T]{values: make(map[T]struct{}, base.GetCap(cap...))}
}

func Collect[T comparable](seq iter.Seq[T]) FrozenSet[T] {
	s := New[T]()
	for t := range seq {
		s.values[t] = struct{}{}
	}
	return s
}

func Consume[S types.MutableSet[T], T comparable](s S) FrozenSet[T] {
	fs := New[T](s.Len())
	for {
		if val, ok := set.Pop(s); ok {
			fs.values[val] = struct{}{}
			continue
		}
		break
	}
	return fs
}

func FromSetLike[S types.SetLike[T], T comparable](s S) FrozenSet[T] {
	switch v := any(s).(type) {
	case set.Set[T]:
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
