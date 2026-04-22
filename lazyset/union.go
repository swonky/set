package lazyset

import "github.com/swonky/set/internal/base"

func NewUnion[S base.SetLike[T], T comparable](sets []S) Union[S, T] {
	return Union[S, T]{sets: append(make([]S, 0, len(sets)), sets...)}
}

type Union[S base.SetLike[T], T any] struct {
	sets []S
}

func (u Union[S, T]) sortByLargest() {
	for i := 1; i < len(u.sets); i++ {
		for j := i; j > 0 && (u.sets[j].Len() > u.sets[j-1].Len()); j-- {
			u.sets[j], u.sets[j-1] = u.sets[j-1], u.sets[j]
		}
	}
}

func (u Union[S, T]) Range(yield func(T) bool) {
	switch len(u.sets) {
	case 0:
		return
	case 1:
		u.sets[0].Range(yield)
		return
	}

	u.sortByLargest()

	u.sets[0].Range(yield)

	for i, s := range u.sets[1:] {
		fn := func(elem T) bool {
			for _, s := range u.sets[:i+1] {
				if s.Contains(elem) {
					return true
				}
			}
			return yield(elem)
		}
		s.Range(fn)
	}
}

func (u Union[S, T]) Len() int {
	n := 0
	u.Range(func(T) bool {
		n++
		return true
	})
	return n
}

func (u Union[S, T]) Contains(elem T) bool {
	for _, s := range u.sets {
		if s.Contains(elem) {
			return true
		}
	}
	return false
}
