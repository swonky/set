package set

import (
	"iter"
	"slices"

	"github.com/swonky/set/guard"
	"github.com/swonky/set/types"
)

func CopyInto[T comparable](dst types.MutableSet[T], src types.SetLike[T]) {
	src.Range(func(t T) bool {
		dst.Add(t)
		return true
	})
}

func CollectInto[T any](dst types.MutableSet[T], src iter.Seq[T]) {
	src(func(t T) bool {
		dst.Add(t)
		return true
	})
}

func SortByLen[T any](sets []types.SetLike[T]) {
	slices.SortFunc(sets, func(a, b types.SetLike[T]) int {
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

func Seq[S types.SetLike[T], T any](s S) iter.Seq[T] {
	return s.Range
}

func Range[S types.SetLike[T], T any](s S, yield func(T) bool) {
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
func GroupBy[T, C comparable](s types.SetLike[T], pred func(T) C) map[C]Set[T] {
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
func Transform[MS types.MutableSet[T], T any](ms MS, fn func(T) T) {
	if fn == nil {
		panic("nil function")
	}
	guard.Write(ms).Do(func(ms2 types.MutableSet[T]) {
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
func TransformInto[S, D comparable](dst types.MutableSet[D], src types.SetLike[S], fn func(S) D) {
	guard.WriteRead(dst, src).Do(
		func(lhs types.MutableSet[D], rhs types.SetLike[S]) {
			yield := func(s S) bool {
				lhs.Add(fn(s))
				return true
			}
			rhs.Range(yield)
		})
}

// All return.
func All[T any](s types.SetLike[T]) iter.Seq2[int, T] {
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
func Values[T any](s types.SetLike[T]) iter.Seq[T] { return s.Range }

// EqualFunc reports whether all elements satisfy fn.
func EqualFunc[T any](s types.SetLike[T], fn func(T) bool) bool {
	result := true

	guard.Read(s).Do(
		func(rset types.SetLike[T]) {
			rset.Range(func(t T) bool {
				if !fn(t) {
					result = false
				}
				return result
			})
		})
	return result
}

// Any reports whether any elements satisfy fn.
func AnyFunc[S types.SetLike[T], T any](s S, fn func(T) bool) bool {
	result := false

	guard.Read(s).Do(
		func(rset types.SetLike[T]) {
			rset.Range(func(t T) bool {
				if fn(t) {
					result = true
					return false
				}
				return true
			})
		})
	return result
}

func IsEqual[S1 types.SetLike[T], S2 types.SetLike[T], T any](s1 S1, s2 S2) bool {
	if guard.Same(s1, s2) {
		return true
	}

	result := true
	guard.ReadRead(s1, s2).Do(
		func(lhs, rhs types.SetLike[T]) {
			if lhs.Len() != rhs.Len() {
				result = false
				return
			}
			lhs.Range(func(t T) bool {
				if !rhs.Contains(t) {
					result = false
				}
				return result
			})
		})
	return result
}

// IsSubset reports whether every element of sub is contained in super.
//
// If either set is lockable, read locks are held for the duration of the operation.
func IsSubset[S1 types.SetLike[T], S2 types.SetLike[T], T any](sub S1, super S2) bool {
	var result bool

	guard.ReadRead(sub, super).Do(
		func(lhs, rhs types.SetLike[T]) {
			if lhs.Len() > rhs.Len() {
				return
			}

			result = true

			lhs.Range(func(t T) bool {
				if !rhs.Contains(t) {
					result = false
					return false
				}
				return true
			})
		})

	return result
}

// AsSlice returns the elements of s as a newly allocated slice.
// The element order is unspecified.
//
// If s is lockable, a read lock is held while the snapshot is collected.
func AsSlice[S types.SetLike[T], T any](s S) []T {
	var result []T

	guard.Read(s).Do(
		func(rset types.SetLike[T]) {
			result = make([]T, 0, rset.Len())

			rset.Range(func(t T) bool {
				result = append(result, t)
				return true
			})
		})

	return result
}

// Filter returns a sequence containing the elements of s for which fn
// returns true.
//
// Filter evaluates fn during iteration. If s is lockable, fn and the
// downstream yield callback run while any read lock held by Range remains
// active.
//
// For synchronized sets, prefer [FilterInto] when materializing
// results under contention or when callbacks may be slow.
func Filter[S types.SetLike[T], T any](s S, fn func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		s.Range(func(t T) bool {
			if fn(t) {
				return yield(t)
			}
			return true
		})
	}
}

// FilterInto adds to dst each element of src for which fn returns true.
//
// FilterInto coordinates access to dst and src using SyncWR. If either set
// is lockable, required locks are held for the duration of the operation.
//
// fn is evaluated while locks are held. It should therefore avoid blocking
// work or calls that may re-enter dst or src.
func FilterInto[T any](dst types.MutableSet[T], src types.SetLike[T], fn func(T) bool) {
	guard.WriteRead(dst, src).Do(
		func(
			lhs types.MutableSet[T],
			rhs types.SetLike[T],
		) {
			rhs.Range(func(t T) bool {
				if fn(t) {
					lhs.Add(t)
				}
				return true
			})
		})
}

// Find returns the first element of s for which fn returns true.
// If no element matches, it returns the zero value of T and false.
//
// If s is lockable, a read lock is held for the duration of the search.
// fn is evaluated while any lock is held.
func Find[T any](s types.SetLike[T], fn func(T) bool) (T, bool) {
	var (
		elem  T
		found bool
	)

	yield := func(t T) bool {
		if fn(t) {
			elem = t
			found = true
			return false
		}
		return true
	}

	guard.Read(s).Do(
		func(rset types.SetLike[T]) {
			rset.Range(yield)
		})

	return elem, found
}

// First returns an arbitrary element of s.
// If s is empty, it returns the zero value of T and false.
//
// If s is lockable, a read lock is held for the duration of the operation.
func First[T any](s types.SetLike[T]) (T, bool) {
	var (
		found bool
		elem  T
	)

	guard.Read(s).Do(
		func(
			rset types.SetLike[T],
		) {
			yield := func(t T) bool {
				elem = t
				found = true
				return false
			}
			rset.Range(yield)
		})

	return elem, found
}

// func op[T any, R any](s types.SetLike[T], fn func(types.SetLike[T]) R) R {
// 	if l, ok := s.(LockableSet[types.MutableSet[T], T]); ok {
// 		var r R
// 		l.WithRLock(func(s2 types.SetLike[T]) {
// 			r = fn(s2)
// 		})
// 		return r
// 	}
// 	return fn(s)
// }

// func WithRead[S types.SetLike[T], T, R any](s S, fn func(types.SetLike[T]) R) R {
// 	if l, ok := types.SetLike[T](s).(LockableSet[types.MutableSet[T], T]); ok {
// 		var r R
// 		l.WithRLock(func(s2 types.SetLike[T]) {
// 			r = fn(s2)
// 		})
// 		return r
// 	}
// 	return fn(s)
// }

// func WithWrite[S types.SetLike[T], T, R any](s S, fn func(S) R) R {
// 	if l, ok := types.SetLike[T](s).(LockableSet[types.MutableSet[T], T]); ok {
// 		var r R
// 		l.WithRLock(func(s2 types.SetLike[T]) { r = fn(s2.(S)) })
// 		return r
// 	}
// 	return fn(s)
// }

// type ReadOperation[S types.SetLike[T], T any, R any] struct {
// 	s S
// 	fn func(types.SetLike[T]) R
// }

// func (v *Operation[S, T, R]) Run() {
// 	if l, ok := v.s.(types.LockableSet[T]); ok {
// 		var r R
// 		l.WithRLock(func(s2 types.SetLike[T]) {
// 			r = fn(s2)
// 		})
// 		return r
// 	}
// 	return fn(s)
// }

// func (v *View) Write[T any, R any](s types.MutableSet[T], fn func(types.MutableSet[T]) R) R {
// 	if l, ok := s.(types.LockableSet[T]); ok {
// 		var r R
// 		l.WithLock(func(s2 types.MutableSet[T]) {
// 			r = fn(s2)
// 		})
// 		return r
// 	}
// 	return fn(s)
// }

// type ReadLocker[S types.SetLike[T], T any] struct {
// 	s S
// }

// type WriteLocker[MS types.MutableSet[T], T any] struct {
// 	ms MS
// }

// func (l ReadLocker[S, T]) WithLock(fn func(types.SetLike[T])) {
// 	if ls, ok := types.SetLike[T](l.s).(LockableSet[types.MutableSet[T], T]); ok {
// 		ls.WithRLock(fn)
// 	}
// }

// func (l WriteLocker[S, T]) WithLock(fn func(types.MutableSet[T])) {
// 	if ls, ok := types.MutableSet[T](l.ms).(LockableSet[types.MutableSet[T], T]); ok {
// 		ls.WithLock(fn)
// 	}
// }

// type Locker[S types.SetLike[T], F types.SetLike[T], T any] interface {
// 	WithLock(func(F))
// }

// var _ Locker[Set[int], SetLike[int], int] = ReadLocker[Set[int], int]{}
// var _ Locker[Set[int], MutableSet[int], int] = WriteLocker[Set[int], int]{}

// type Pair2[S1, F1 SetLike[T1], S2, F2 SetLike[T2], T1, T2 any] struct {
// 	lhs Locker[S1, F1, T1]
// 	rhs Locker[S2, F2, T2]
// }

// func (p *Pair2[S1, F1, S2, F2, T1, T2]) WithLock(fn func(SetLike[T1], SetLike[T2])) {
// 	p.lhs.WithLock(fn)
// }

// type Pair[T1 any, T2 any] struct {
// 	lhs SetLike[T1]
// 	rhs SetLike[T2]
// }

// func (p *Pair[T1, T2]) WithReadLock(fn func(SetLike[T1], SetLike[T2])) {
// 	if la, ok := p.lhs.(LockableSet[MutableSet[T1], T1]); ok {
// 		if lb, ok := p.rhs.(LockableSet[MutableSet[T2], T2]); ok {
// 			la.WithRLock(func(a2 SetLike[T1]) {
// 				lb.WithRLock(func(b2 SetLike[T2]) {
// 					fn(a2, b2)
// 				})
// 			})
// 		}
// 	}
// 	fn(p.lhs, p.rhs)
// }

// func NewPair[T1 any, T2 any](lhs SetLike[T1], rhs SetLike[T2]) *Pair[T1, T2] {
// 	return &Pair[T1, T2]{lhs: lhs, rhs: rhs}
// }

// func op2[A any, B any, R any](
// 	a SetLike[A],
// 	b SetLike[B],
// 	fn func(a SetLike[A], b SetLike[B]) R,
// ) R {
// 	if la, ok := a.(LockableSet[MutableSet[A], A]); ok {
// 		if lb, ok := b.(LockableSet[MutableSet[B], B]); ok {
// 			var r R
// 			la.WithRLock(func(a2 SetLike[A]) {
// 				lb.WithRLock(func(b2 SetLike[B]) {
// 					r = fn(a2, b2)
// 				})
// 			})
// 			return r
// 		}
// 	}
// 	return fn(a, b)
// }
