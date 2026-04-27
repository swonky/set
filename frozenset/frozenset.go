package frozenset

import (
	"maps"

	"github.com/swonky/set"
	"github.com/swonky/set/types"
)

var _ types.SetLike[int] = FrozenSet[int]{}

type FrozenSet[T comparable] struct {
	values map[T]struct{}
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

func (fs FrozenSet[T]) AsSet() set.Set[T] {
	return maps.Clone(fs.values)
}
