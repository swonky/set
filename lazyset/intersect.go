// intersect.go
package lazyset

import (
	"iter"

	"github.com/swonky/set/internal/base"
)

var _ base.SetLike[int] = Intersection[base.SetLike[int], int]{}

func NewIntersection[S base.SetLike[T], T comparable](sets []S) Intersection[S, T] {
	return Intersection[S, T]{sets: append(make([]S, 0, len(sets)), sets...)}
}

func (li Intersection[S, T]) sortBySmallest() {
	for i := 1; i < len(li.sets); i++ {
		for j := i; j > 0 && (li.sets[j].Len() < li.sets[j-1].Len()); j-- {
			li.sets[j], li.sets[j-1] = li.sets[j-1], li.sets[j]
		}
	}
}

type Intersection[S base.SetLike[T], T any] struct {
	sets []S
}

func (li Intersection[S, T]) Sets() iter.Seq[S] {
	return func(yield func(S) bool) {
		for _, s := range li.sets {
			if !yield(s) {
				return
			}
		}
	}
}

func (li Intersection[S, T]) Range(yield func(T) bool) {
	switch len(li.sets) {
	case 0:
		return
	case 1:
		li.sets[0].Range(yield)
		return
	}

	li.sortBySmallest()

	fn := func(elem T) bool {
		for _, s := range li.sets[1:] {
			if !s.Contains(elem) {
				return true
			}
		}
		return yield(elem)
	}

	smallest := li.sets[0]

	if ls, ok := any(smallest).(base.LockableSet[T]); ok {
		ls.WithRLock(func(s base.SetLike[T]) {
			s.Range(fn)
		})
		return
	}
	smallest.Range(fn)
}

func (li Intersection[S, T]) Contains(elem T) bool {
	for _, s := range li.sets {
		if !s.Contains(elem) {
			return false
		}
	}
	return true
}

func (li Intersection[S, T]) Len() int {
	n := 0
	li.Range(func(T) bool {
		n++
		return true
	})
	return n
}
