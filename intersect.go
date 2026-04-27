// intersect.go
package set

import (
	"iter"

	"github.com/swonky/set/types"
)

var _ types.SetLike[int] = Intersection[int]{}

// func Intersect[S types.SetLike[T], T comparable](sets []S) Intersection[T] {
func Intersect[S types.SetLike[T], T comparable](sets ...S) Intersection[T] {
	return Intersection[T]{sets: append(make([]types.SetLike[T], 0, len(sets)), any(sets).([]types.SetLike[T])...)}
}

func (is Intersection[T]) sortBySmallest() {
	for i := 1; i < len(is.sets); i++ {
		for j := i; j > 0 && (is.sets[j].Len() < is.sets[j-1].Len()); j-- {
			is.sets[j], is.sets[j-1] = is.sets[j-1], is.sets[j]
		}
	}
}

type Intersection[T any] struct {
	sets []types.SetLike[T]
}

func (is Intersection[T]) Sets() iter.Seq[types.SetLike[T]] {
	return func(yield func(types.SetLike[T]) bool) {
		for _, s := range is.sets {
			if !yield(s) {
				return
			}
		}
	}
}

func (is Intersection[T]) Range(yield func(T) bool) {
	switch len(is.sets) {
	case 0:
		return
	case 1:
		is.sets[0].Range(yield)
		return
	}

	is.sortBySmallest()

	fn := func(elem T) bool {
		for _, s := range is.sets[1:] {
			if !s.Contains(elem) {
				return true
			}
		}
		return yield(elem)
	}

	smallest := is.sets[0]

	if ls, ok := any(smallest).(types.LockableSet[T]); ok {
		ls.WithRLock(func(s types.SetLike[T]) {
			s.Range(fn)
		})
		return
	}
	smallest.Range(fn)
}

func (is Intersection[T]) Contains(elem T) bool {
	for _, s := range is.sets {
		if !s.Contains(elem) {
			return false
		}
	}
	return true
}

func (is Intersection[T]) Len() int {
	n := 0
	is.Range(func(T) bool {
		n++
		return true
	})
	return n
}
