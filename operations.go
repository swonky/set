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
	WithWrite(ms, func(ms2 MS) struct{} {
		for elem := range Consume(ms2) {
			ms2.Add(fn(elem))
		}
		return struct{}{}
	})
}

// TransformInto maps fn to each element of src and adds the result to dst.
//
// TransformInto panics if fn is nil.
func TransformInto[S, D comparable](dst MutableSet[D], src SetLike[S], fn func(S) D) {
	op2(src, dst,
		func(a SetLike[S], b SetLike[D]) struct{} {
			mb := b.(MutableSet[D])
			a.Range(
				func(s S) bool {
					mb.Add(fn(s))
					return true
				},
			)
			return struct{}{}
		},
	)
}

// All reports whether all elements satisfy fn.
func All[T any](s SetLike[T]) iter.Seq[T] { return s.Range }

// EqualFunc reports whether all elements satisfy fn.
func EqualFunc[T any](s SetLike[T], fn func(T) bool) bool {
	return op(s, func(sl SetLike[T]) bool {
		for t := range sl.Range {
			if !fn(t) {
				return false
			}
		}
		return true
	})
}

// Any reports whether any elements satisfy fn.
func AnyFunc[T any](s SetLike[T], fn func(T) bool) bool {
	return op(s, func(sl SetLike[T]) bool {
		for t := range sl.Range {
			if fn(t) {
				return true
			}
		}
		return false
	})
}

func IsEqual[T any](a, b SetLike[T]) bool {
	return op2(a, b, func(a, b SetLike[T]) bool {
		if a.Len() != b.Len() {
			return false
		}
		for k := range a.Range {
			if !b.Contains(k) {
				return false
			}
		}
		return true
	})
}

func IsSubset[T any](a, b SetLike[T]) bool {
	return op2(a, b, func(a, b SetLike[T]) bool {
		if a.Len() > b.Len() {
			return false
		}
		for k := range a.Range {
			if b.Contains(k) {
				return false
			}
		}
		return true
	})
}

func AsSlice[T any](s SetLike[T]) []T {
	return op(s, func(sl SetLike[T]) []T {
		out := make([]T, 0, s.Len())
		for k := range s.Range {
			out = append(out, k)
		}
		return out
	})
}

func Filter[T any](s SetLike[T], fn func(T) bool) iter.Seq[T] {
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
