package lazyset

import (
	"github.com/swonky/set/internal/base"
)

var _ base.SetLike[int] = LazySet[int, BinaryOp[int]]{}

type BinaryOp[T any] interface {
	Range(a, b base.SetLike[T], yield func(T) bool)
	Contains(a, b base.SetLike[T], elem T) bool
}

type LazySet[T comparable, B BinaryOp[T]] struct {
	a, b base.SetLike[T]
	op   B
}

func New[T comparable, B BinaryOp[T]](
	sets []base.SetLike[T],
	op B,
) *LazySet[T, B] {
	if len(sets) == 1 {
		return &LazySet[T, B]{a: sets[0], b: sets[0], op: op}
	}
	acc := &LazySet[T, B]{a: sets[0], b: sets[1], op: op}
	for _, s := range sets[2:] {
		acc = &LazySet[T, B]{a: acc, b: s, op: op}
	}
	return acc
}

func (s LazySet[T, B]) Len() int {
	n := 0
	s.Range(func(T) bool {
		n++
		return true
	})
	return n
}

func (s LazySet[T, B]) AsSet() base.Set[T] {
	bs := make(base.Set[T])
	s.Range(func(t T) bool {
		bs[t] = struct{}{}
		return true
	})
	return bs
}

func (s LazySet[T, B]) AsSet2() base.Set[T] {
	bs := make(base.Set[T], s.Len())
	s.Range(func(t T) bool {
		bs[t] = struct{}{}
		return true
	})
	return bs
}

type Data struct {
	a, b, c int
}

func StackAlloc() Data {
	return Data{1, 2, 3} // stays on stack
}

func HeapAlloc() *Data {
	return &Data{1, 2, 3} // escapes to heap
}

func (s LazySet[T, B]) Range(yield func(T) bool) {
	if la, ok := s.a.(base.LockableSet[T]); ok {
		if lb, ok := s.b.(base.LockableSet[T]); ok {
			la.WithRLock(func(a2 base.SetLike[T]) {
				lb.WithRLock(func(b2 base.SetLike[T]) {
					s.op.Range(a2, b2, yield)
				})
			})
			return
		}
	}
	s.op.Range(s.a, s.b, yield)
}

func (s LazySet[T, B]) Contains(elem T) bool {
	return s.op.Contains(s.a, s.b, elem)
}
