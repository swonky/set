package set

import (
	"iter"
	"slices"
)

func CopyInto[T comparable](dst MutableSet[T], src SetLike[T]) {
	src.Range(func(t T) bool {
		dst.Add(t)
		return true
	})
}

func CollectInto[T any](dst MutableSet[T], src iter.Seq[T]) {
	src(func(t T) bool {
		dst.Add(t)
		return true
	})
}

func SortByLen[T any](sets []SetLike[T]) {
	slices.SortFunc(sets, func(a, b SetLike[T]) int {
		la, lb := a.Len(), b.Len()
		switch {
		case la < lb:
			return -1
		case la > lb:
			return 1
		default:
			return 0
		}
	})
}

func Seq[S SetLike[T], T any](s S) iter.Seq[T] {
	return s.Range
}

func Range[S SetLike[T], T any](s S, yield func(T) bool) {
	s.Range(yield)
}

// GroupBy groups the elements of s by the value returned from pred.
//
// Each element k in s is assigned to a group keyed by pred(k). The returned map
// contains one entry per distinct key, where each value is a Set[T] of all
// elements that produced that key.
//
// The grouping is exhaustive and disjoint: every element appears in exactly one
// group. Iteration order is not preserved.
//
// GroupBy panics if pred is nil.
func GroupBy[T, C comparable](s SetLike[T], pred func(T) C) map[C]Set[T] {
	if pred == nil {
		panic("nil predicate in set.GroupBy")
	}
	r := make(map[C]Set[T])
	for k := range s.Range {
		c := pred(k)
		g, ok := r[c]
		if !ok {
			g = make(Set[T])
			r[c] = g
		}
		g[k] = struct{}{}
	}
	return r
}

// Transform replaces each element in ms with fn(elem).
// If multiple elements produce the same result, duplicates are discarded.
//
// For lockable sets, Transform holds an exclusive lock for the duration
// of the operation.
//
// fn must not be nil and must not call methods on ms.
// If fn panics, ms may be left partially transformed.
func Transform[MS MutableSet[T], T any](ms MS, fn func(T) T) {
	if fn == nil {
		panic("nil function")
	}
	Sync(ms).W(func(ms2 MutableSet[T]) {
		Consume(ms2)(
			func(t T) bool {
				ms2.Add(fn(t))
				return true
			})
	})
}

// TransformInto maps fn to each element of src and adds the result to dst.
//
// TransformInto panics if fn is nil.
func TransformInto[S, D comparable](dst MutableSet[D], src SetLike[S], fn func(S) D) {
	Sync2(dst, src).WR(
		func(ms MutableSet[D], s SetLike[S]) {
			s.Range(
				func(s S) bool {
					ms.Add(fn(s))
					return true
				})
		})

}

// All return.
func All[T any](s SetLike[T]) iter.Seq2[int, T] {
	i := 0
	return func(yield func(int, T) bool) {
		s.Range(func(t T) bool {
			r := yield(i, t)
			i++
			return r
		})
	}
}

// All reports whether all elements satisfy fn.
func Values[T any](s SetLike[T]) iter.Seq[T] { return s.Range }

// EqualFunc reports whether all elements satisfy fn.
func EqualFunc[T any](s SetLike[T], fn func(T) bool) bool {
	return op(s,
		func(sl SetLike[T]) bool {
			for t := range sl.Range {
				if !fn(t) {
					return false
				}
			}
			return true
		})
}

// Any reports whether any elements satisfy fn.
func AnyFunc[S SetLike[T], T any](s S, fn func(T) bool) bool {
	return op(s,
		func(sl SetLike[T]) bool {
			for t := range sl.Range {
				if fn(t) {
					return true
				}
			}
			return false
		})
}

func IsEqual[S1 SetLike[T], S2 SetLike[T], T any](s1 S1, s2 S2) bool {
	var result bool

	NewPair(s1, s2).WithReadLock(
		func(sl1, sl2 SetLike[T]) {
			if sl1.Len() != sl2.Len() {
				return
			}
			for k := range sl1.Range {
				if !sl2.Contains(k) {
					return
				}
			}
			result = true
		})

	return result

	// return op2(s1, s2,
	// 	func(a, b SetLike[T]) bool {
	// 		if a.Len() != b.Len() {
	// 			return false
	// 		}
	// 		for k := range a.Range {
	// 			if !b.Contains(k) {
	// 				return false
	// 			}
	// 		}
	// 		return true
	// 	})
}

func IsSubset[S1 SetLike[T], S2 SetLike[T], T any](sub S1, super S2) bool {
	return op2(sub, super,
		func(sub2, super2 SetLike[T]) bool {
			if sub2.Len() > super2.Len() {
				return false
			}
			for k := range Seq(sub2) {
				if super2.Contains(k) {
					return false
				}
			}
			return true
		})
}

func AsSlice[S SetLike[T], T any](s S) []T {
	return WithRead(s,
		func(sl SetLike[T]) []T {
			out := make([]T, 0, s.Len())
			for k := range s.Range {
				out = append(out, k)
			}
			return out
		})
}

func Filter[S SetLike[T], T any](s S, fn func(T) bool) iter.Seq[T] {
	return op(s, func(sl SetLike[T]) iter.Seq[T] {
		return func(yield func(T) bool) {
			sl.Range(func(t T) bool {
				if fn(t) {
					return yield(t)
				}
				return true
			})
		}
	})
}

func FilterInto[T any](dst MutableSet[T], src SetLike[T], fn func(T) bool) {
	op2(src, dst, func(src, dst SetLike[T]) struct{} {
		ms := dst.(MutableSet[T])
		for t := range src.Range {
			if fn(t) {
				ms.Add(t)
			}
		}
		return struct{}{}
	})
}

func Find[T any](s SetLike[T], fn func(T) bool) (T, bool) {
	type result[T any] struct {
		v  T
		ok bool
	}

	r := op(s, func(sl SetLike[T]) result[T] {
		var out result[T]
		sl.Range(func(t T) bool {
			if fn(t) {
				out.v = t
				out.ok = true
				return false
			}
			return true
		})
		return out
	})
	return r.v, r.ok
}

func First[T any](s SetLike[T]) (T, bool) {
	type result[T any] struct {
		v  T
		ok bool
	}
	r := op(s, func(sl SetLike[T]) result[T] {
		var out result[T]
		sl.Range(func(t T) bool {
			out.v = t
			out.ok = true
			return false
		})
		return out
	})
	return r.v, r.ok
}

func op[T any, R any](s SetLike[T], fn func(SetLike[T]) R) R {
	if l, ok := s.(LockableSet[T]); ok {
		var r R
		l.WithRLock(func(s2 SetLike[T]) {
			r = fn(s2)
		})
		return r
	}
	return fn(s)
}

func WithRead[S SetLike[T], T, R any](s S, fn func(SetLike[T]) R) R {
	if l, ok := SetLike[T](s).(LockableSet[T]); ok {
		var r R
		l.WithRLock(func(s2 SetLike[T]) {
			r = fn(s2)
		})
		return r
	}
	return fn(s)
}

func WithWrite[S SetLike[T], T, R any](s S, fn func(S) R) R {
	if l, ok := SetLike[T](s).(LockableSet[T]); ok {
		var r R
		l.WithRLock(func(s2 SetLike[T]) { r = fn(s2.(S)) })
		return r
	}
	return fn(s)
}

// type ReadOperation[S SetLike[T], T any, R any] struct {
// 	s S
// 	fn func(SetLike[T]) R
// }

// func (v *Operation[S, T, R]) Run() {
// 	if l, ok := v.s.(LockableSet[T]); ok {
// 		var r R
// 		l.WithRLock(func(s2 SetLike[T]) {
// 			r = fn(s2)
// 		})
// 		return r
// 	}
// 	return fn(s)
// }

// func (v *View) Write[T any, R any](s MutableSet[T], fn func(MutableSet[T]) R) R {
// 	if l, ok := s.(LockableSet[T]); ok {
// 		var r R
// 		l.WithLock(func(s2 MutableSet[T]) {
// 			r = fn(s2)
// 		})
// 		return r
// 	}
// 	return fn(s)
// }

type ReadLocker[S SetLike[T], T any] struct {
	s S
}

type WriteLocker[MS MutableSet[T], T any] struct {
	ms MS
}

func (l ReadLocker[S, T]) WithLock(fn func(SetLike[T])) {
	if ls, ok := SetLike[T](l.s).(LockableSet[T]); ok {
		ls.WithRLock(fn)
	}
}

func (l WriteLocker[S, T]) WithLock(fn func(MutableSet[T])) {
	if ls, ok := MutableSet[T](l.ms).(LockableSet[T]); ok {
		ls.WithLock(fn)
	}
}

type Locker[S SetLike[T], F SetLike[T], T any] interface {
	WithLock(func(F))
}

var _ Locker[Set[int], SetLike[int], int] = ReadLocker[Set[int], int]{}
var _ Locker[Set[int], MutableSet[int], int] = WriteLocker[Set[int], int]{}

type Pair2[S1, F1 SetLike[T1], S2, F2 SetLike[T2], T1, T2 any] struct {
	lhs Locker[S1, F1, T1]
	rhs Locker[S2, F2, T2]
}

func (p *Pair2[S1, F1, S2, F2, T1, T2]) WithLock(fn func(SetLike[T1], SetLike[T2])) {
	p.lhs.WithLock()
}

type Pair[T1 any, T2 any] struct {
	lhs SetLike[T1]
	rhs SetLike[T2]
}

func (p *Pair[T1, T2]) WithReadLock(fn func(SetLike[T1], SetLike[T2])) {
	if la, ok := p.lhs.(LockableSet[T1]); ok {
		if lb, ok := p.rhs.(LockableSet[T2]); ok {
			la.WithRLock(func(a2 SetLike[T1]) {
				lb.WithRLock(func(b2 SetLike[T2]) {
					fn(a2, b2)
				})
			})
		}
	}
	fn(p.lhs, p.rhs)
}

func NewPair[T1 any, T2 any](lhs SetLike[T1], rhs SetLike[T2]) *Pair[T1, T2] {
	return &Pair[T1, T2]{lhs: lhs, rhs: rhs}
}

func op2[A any, B any, R any](
	a SetLike[A],
	b SetLike[B],
	fn func(a SetLike[A], b SetLike[B]) R,
) R {
	if la, ok := a.(LockableSet[A]); ok {
		if lb, ok := b.(LockableSet[B]); ok {
			var r R
			la.WithRLock(func(a2 SetLike[A]) {
				lb.WithRLock(func(b2 SetLike[B]) {
					r = fn(a2, b2)
				})
			})
			return r
		}
	}
	return fn(a, b)
}
