package set

import (
	"iter"
	"maps"
	"slices"

	"github.com/swonky/set/internal/base"
	"github.com/swonky/set/lazyset"
)

// New returns a new Set with a .
// If no items are provided, it returns an empty set.
func New[T comparable](cap ...int) Set[T] {
	return make(Set[T], base.GetCap(cap...))
}

// Collect returns a new Set containing all elements produced by the iterator.
// If a positive capacity hint is provided, it is used to preallocate the set.
func Collect[T comparable](it iter.Seq[T], size ...int) Set[T] {
	s := make(base.Set[T], base.GetCap(size...))
	for v := range it {
		s[v] = struct{}{}
	}
	return s
}

func FromSetLike[T comparable](s SetLike[T]) Set[T] {
	switch e := s.(type) {
	case Set[T]:
		return maps.Clone(e)
	case AsSetter[T]:
		return e.AsSet()
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

// Specific reducer functions

// // Union returns a new Set containing all elements from the provided sets.
// func Union[T comparable](sets ...SetLike[T]) Set[T] {
// 	return Reduce(sets, func(a, b Set[T]) Set[T] { return a.Union(b) })
// }

// func unionImpl[T comparable](s []Set[T]) Set[T] {
// 	return Reduce(s, func(a, b Set[T]) Set[T] { return a.Union(b) })
// }

// Intersect returns a new Set containing elements present in all input sets.
// Evaluation stops early if the result becomes empty.
// func Intersect[T comparable](sets ...Set[T]) Set[T] {
// 	return ReduceTry(
// 		sets,
// 		func(a, b Set[T]) (Set[T], bool) {
// 			r := a.Intersect(b)
// 			return r, !r.IsEmpty()
// 		},
// 	)
// }

// func Union[T comparable](sets []SetLike[T]) *lazyset.LazySet[T, UnionOp[T]] {
// 	slices.SortFunc(sets, func(a, b SetLike[T]) int { return cmp.Compare(b.Len(), a.Len()) })
// 	return lazyset.New(sets, UnionOp[T]{})
// }

type UnionOp[T any] struct{}

func (UnionOp[T]) Range(a, b base.SetLike[T], yield func(T) bool) {
	a.Range(yield)
	b.Range(func(elem T) bool {
		if a.Contains(elem) {
			return true
		}
		return yield(elem)
	})
}

func (UnionOp[T]) Contains(a, b base.SetLike[T], elem T) bool {
	return a.Contains(elem) || b.Contains(elem)
}

func unionContainsFn[T any](a, b base.SetLike[T], elem T) bool {
	return a.Contains(elem) || b.Contains(elem)
}

func intersectRangeFn[T any](a, b base.SetLike[T], yield func(T) bool) {
	a.Range(func(elem T) bool {
		if !b.Contains(elem) {
			return true
		}
		return yield(elem)
	})
}

func intersectContainsFn[T any](a, b base.SetLike[T], elem T) bool {
	return a.Contains(elem) && b.Contains(elem)
}

// ops.go

// func Union[T any](sets ...SetLike[T]) lazyset.LazySet[T] {
// 	return lazyset.New(unionTwo[T], sets)
// }

// // unionTwo yields elements from b that are not already present in a.
// // a represents all sets yielded so far, so this produces no duplicates
// // across the full reduction.
// func unionTwo[T any](a, b base.SetLike[T], yield func(T) bool) bool {
// 	if la, ok := a.(base.LockableSet[T]); ok {
// 		var cont bool
// 		la.WithRLock(func(ua base.SetLike[T]) {
// 			cont = unionTwo(ua, b, yield)
// 		})
// 		return cont
// 	}
// 	if lb, ok := b.(base.LockableSet[T]); ok {
// 		var cont bool
// 		lb.WithRLock(func(ub base.SetLike[T]) {
// 			cont = unionTwo(a, ub, yield)
// 		})
// 		return cont
// 	}
// 	b.Range(func(elem T) bool {
// 		if a.Contains(elem) {
// 			return true
// 		}
// 		return yield(elem)
// 	})
// 	return true
// }

// func Union[T any](sets ...SetLike[T]) lazyset.LazySet[T] {
// 	return lazyset.New(unionImpl[T], sets)
// }

func UnionInto[T any](out MutableSet[T], sets ...SetLike[T]) {
	// ordered := slices.Clone(sets)
	// SortByLen(ordered)
	for _, s := range sets {
		for elem := range s.Range {
			out.Add(elem)
		}
	}
}

func CopyInto[T comparable](out Set[T], sets ...Set[T]) {
	// ordered := slices.Clone(sets)
	// SortByLen(ordered)
	for _, s := range sets {
		maps.Copy(out, s)
	}
}

// func unionImpl[T any](a, b SetLike[T]) {
// 	for elem := range UnionIter(ordered...) {
// 		if len(a) > len(b) {
// 			a, b = b, a
// 		}
// 		r := s.Clone()
// 		for k := range o {
// 			r[k] = struct{}{}
// 		}
// 		return r
// 	}
// }

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

// func IntersectIter[T comparable](sets []SetLike[T]) iter.Seq[T] {
// 	return func(yield func(T) bool) { IntersectRange(sets, yield) }
// }

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

// TransformIter returns an iterator that yields fn(k) for each element k in s.
//
// The transformation is applied lazily during iteration and does not allocate
// an intermediate collection.
//
// Unlike Transform, duplicate results are not removed. If multiple elements map
// to the same value, duplicates will be yielded. Iteration order is not
// guaranteed.
//
// TransformIter panics if fn is nil.
func TransformIter[T comparable, U any](s SetLike[T], fn func(T) U) iter.Seq[U] {
	if fn == nil {
		panic("nil function")
	}
	return func(yield func(U) bool) {
		for k := range s.Range {
			if !yield(fn(k)) {
				return
			}
		}
	}
}

func Intersect[T comparable](sets ...SetLike[T]) Intersection[SetLike[T], T] {
	return lazyset.NewIntersection(sets)
}

func Unite[T comparable](sets ...SetLike[T]) Union[SetLike[T], T] {
	return lazyset.NewUnion(sets)
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

// Transform applies fn to each element of s and returns a new Set containing
// the results.
//
// Each element k in s is transformed to fn(k). Because the result is a set,
// duplicate outputs are deduplicated.
//
// The operation is eager and allocates a new set. Iteration order is not
// preserved.
//
// Transform panics if fn is nil.
func Transform[A, B comparable](s SetLike[A], fn func(A) B) Set[B] {
	return op(s, func(sl SetLike[A]) Set[B] {
		r := New[B](s.Len())
		for t := range sl.Range {
			r.Add(fn(t))
		}
		return r
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

func FilterSet[T comparable](s SetLike[T], fn func(T) bool) Set[T] {
	return op(s, func(sl SetLike[T]) Set[T] {
		r := New[T](sl.Len())
		for t := range sl.Range {
			if fn(t) {
				r.Add(t)
			}
		}
		return r
	})
}

func FilterInto[T any](in SetLike[T], out MutableSet[T], fn func(T) bool) {
	op2(in, out, func(src, dst SetLike[T]) struct{} {
		ms := dst.(MutableSet[T])
		for t := range src.Range {
			if fn(t) {
				ms.Add(t)
			}
		}
		return struct{}{}
	})
}

func filterValues[K comparable, V any, M ~map[K]V](m M, fn func(V) bool) M {
	out := make(M, len(m))
	for k, v := range m {
		if fn(v) {
			out[k] = v
		}
	}
	return out
}

type result[T any] struct {
	v  T
	ok bool
}

func Find[T any](s SetLike[T], fn func(T) bool) (T, bool) {
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

// AddCheck inserts an item into the set. Returns true if item was already present.
func AddCheck[T any](ms MutableSet[T], elem T) bool {
	if ms.Contains(elem) {
		return true
	}
	ms.Add(elem)
	return false
}

func Append[T any](ms MutableSet[T], elems ...T) {
	for _, v := range elems {
		ms.Add(v)
	}
}

func Remove[T any](ms MutableSet[T], elems ...T) {
	for _, v := range elems {
		ms.Delete(v)
	}
}

func Pop[T any](ms MutableSet[T]) (T, bool) {
	out := op(ms, func(sl SetLike[T]) result[T] {
		msl := sl.(MutableSet[T])
		var r result[T]
		for t := range sl.Range {
			r.v = t
			r.ok = true
			msl.Delete(t)
			return r
		}
		return r
	})
	return out.v, out.ok
}

func mapPop[K comparable, V any](m map[K]V) (K, V, bool) {
	for k, v := range m {
		delete(m, k)
		return k, v, true
	}
	var k K
	var v V
	return k, v, false
}

func Clear[T any](s MutableSet[T]) {
	if s.Len() == 0 {
		return
	}
	op(s, func(sl SetLike[T]) struct{} {
		ms := s
		sl.Range(func(t T) bool {
			ms.Delete(t)
			return true
		})
		return struct{}{}
	})
}

func UnionWith[T any](dst MutableSet[T], src SetLike[T]) {
	op2(src, dst, func(a, b SetLike[T]) struct{} {
		ms := b.(MutableSet[T])

		a.Range(func(t T) bool {
			ms.Add(t)
			return true
		})
		return struct{}{}
	})
}

func IntersectWith[T any](dst MutableSet[T], src SetLike[T]) {
	op2(dst, src, func(a, b SetLike[T]) struct{} {
		ms := a.(MutableSet[T])

		a.Range(func(t T) bool {
			if !b.Contains(t) {
				ms.Delete(t)
			}
			return true
		})
		return struct{}{}
	})
}

func DiffWith[T any](dst MutableSet[T], src SetLike[T]) {
	op2(dst, src, func(a, b SetLike[T]) struct{} {
		md := a.(MutableSet[T])
		for t := range b.Range {
			if md.Contains(t) {
				dst.Delete(t)
			}
		}
		return struct{}{}
	})
}

func SymDiffWith[T any](dst MutableSet[T], src SetLike[T]) {
	op2(dst, src, func(a, b SetLike[T]) struct{} {
		md := a.(MutableSet[T])
		for t := range b.Range {
			if md.Contains(t) {
				dst.Delete(t)
			} else {
				dst.Add(t)
			}
		}
		return struct{}{}
	})
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

func op2[T any, R any](
	a, b SetLike[T],
	fn func(a, b SetLike[T]) R,
) R {
	if la, ok := a.(LockableSet[T]); ok {
		if lb, ok := b.(LockableSet[T]); ok {
			var r R
			la.WithRLock(func(a2 SetLike[T]) {
				lb.WithRLock(func(b2 SetLike[T]) {
					r = fn(a2, b2)
				})
			})
			return r
		}
	}
	return fn(a, b)
}

func CollectInto[T any](dst MutableSet[T], src iter.Seq[T]) {
	for t := range src {
		dst.Add(t)
	}
}
