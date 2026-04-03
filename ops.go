package set

import (
	"iter"
	"slices"
)

var _ SetLike[int] = FrozenSet[int]{}

// New returns a new Set containing the provided items.
// If no items are provided, it returns an empty set.
func New[T comparable](items ...T) Set[T] {
	s := make(Set[T], len(items))
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

// Make returns a new Set with a .
// If no items are provided, it returns an empty set.
func Make[T comparable](size ...int) Set[T] {
	if len(size) > 0 && size[0] > 0 {
		return make(Set[T], size[0])
	}
	return make(Set[T], 0)
}

// Collect returns a new Set containing all elements produced by the iterator.
// If a positive capacity hint is provided, it is used to preallocate the set.
func Collect[T comparable](it iter.Seq[T], size ...int) Set[T] {
	s := Make[T](size...)
	for v := range it {
		s[v] = struct{}{}
	}
	return s
}

func Clone[T any](s Cloner[T]) T {
	return s.Clone()
}

// Specific reducer functions

// Union returns a new Set containing all elements from the provided sets.
func Union[T comparable](sets ...Set[T]) Set[T] {
	return Reduce(sets, func(a, b Set[T]) Set[T] { return a.Union(b) })
}

// Intersect returns a new Set containing elements present in all input sets.
// Evaluation stops early if the result becomes empty.
func Intersect[T comparable](sets ...Set[T]) Set[T] {
	return ReduceTry(
		sets,
		func(a, b Set[T]) (Set[T], bool) {
			r := a.Intersect(b)
			return r, !r.IsEmpty()
		},
	)
}

// UnionIter returns an iterator of unique elements from the provided sets,
// preserving first occurrence order across sets.
func UnionIter[T comparable](sets ...SetLike[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := New[T]()
		for _, s := range sets {
			for e := range s.Range {
				if _, exists := seen[e]; !exists {
					seen[e] = struct{}{}
					if !yield(e) {
						return
					}
				}
			}
		}
	}
}

func SortByLen[T comparable](sets []SetLike[T]) {
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

// IntersectRange calls fn for each element that exists in all provided sets.
//
// It is intended for high-frequency filtering across multiple membership sets where early exit matters.
// It performs no allocations and does not construct a result set.
//
// For each matching element k, fn(k) is invoked.
// If fn returns false, iteration stops immediately.
// The order of elements is unspecified.
//
// IntersectRange must not be used concurrently with modifications to any of the input sets.
// The behavior is undefined if the sets are changed while iteration is in progress.
// The function fn must not modify the input sets.
//
// Use IntersectRange when you only need to test or process elements on the fly.
// Use [Intersect] when you need to retain the full result.
//
// Example: apply an operation to all elements present in every set,
// without allocating an intermediate result.
//
//	IntersectRange([]SetLike[int]{a, b, c}, func(x int) bool {
//		process(x)
//		return true
//	})
//
// Example: check whether a value exists in all sets
//
//	found := false
//	IntersectRange([]SetLike[string]{setA, setB, setC}, func(s string) bool {
//		if s == "foo" {
//			found = true
//			return false
//		}
//		return true
//	})
// func IntersectRange[T comparable](
// 	sets []SetLike[T],
// 	fn func(T) bool,
// ) {
// 	switch len(sets) {
// 	case 0:
// 		return
// 	case 1:
// 		sets[0].Range(fn)
// 		return
// 	}

// 	smallest := 0
// 	for i := 1; i < len(sets); i++ {
// 		if sets[i].Len() < sets[smallest].Len() {
// 			smallest = i
// 		}
// 	}

// 	s0 := sets[smallest]

// 	switch s := any(s0).(type) {
// 	case Set[T]:
// 		for k := range s {
// 			if !intersectHasAll(sets, smallest, k) {
// 				continue
// 			}
// 			if !fn(k) {
// 				return
// 			}
// 		}
// 		return
// 	case FrozenSet[T]:
// 		for k := range s.s {
// 			if !intersectHasAll(sets, smallest, k) {
// 				continue
// 			}
// 			if !fn(k) {
// 				return
// 			}
// 		}
// 		return

// 	case *SyncSet[T]:
// 		s.WithReadLock(func(sl SetLike[T]) {
// 			for k := range sl.Range {
// 				if !intersectHasAll(sets, smallest, k) {
// 					continue
// 				}
// 				if !fn(k) {
// 					return
// 				}
// 			}
// 			return
// 		})

// 	}

// 	s0.Range(func(k T) bool {
// 		if !intersectHasAll(sets, smallest, k) {
// 			return true
// 		}
// 		return fn(k)
// 	})
// }

func IntersectRangeNew[T comparable](
	sets []SetLike[T],
	fn func(T) bool,
) {
	switch len(sets) {
	case 0:
		return
	case 1:
		sets[0].Range(fn)
		return
	}

	s0, smallest := getSmallestSet(sets)

	switch s := any(s0).(type) {
	case Set[T]:
		s.Range(itersectRanger(sets, smallest, fn))
	case FrozenSet[T]:
		s.Range(itersectRanger(sets, smallest, fn))
	case *SyncSet[T]:
		s.Range(itersectRanger(sets, smallest, fn))
	default:
		s0.Range(itersectRanger(sets, smallest, fn))
	}
}

func itersectRanger[T comparable](sets []SetLike[T], skip int, fn func(t T) bool) func(T) bool {
	return func(t T) bool {
		for i := range sets {
			if i == skip {
				continue
			}
			if !sets[i].Contains(t) {
				return false
			}
		}
		return fn(t)
	}
}

func intersectHasAll[T comparable](sets []SetLike[T], skip int, k T) bool {
	for i := range sets {
		if i == skip {
			continue
		}
		if !sets[i].Contains(k) {
			return false
		}
	}
	return true
}

func getSmallestSet[T comparable](sets []SetLike[T]) (SetLike[T], int) {
	smallest := 0
	for i := 1; i < len(sets); i++ {
		if sets[i].Len() < sets[smallest].Len() {
			smallest = i
		}
	}
	return sets[smallest], smallest
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
		r := Make[B](s.Len())
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
		r := Make[T](sl.Len())
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
		l.WithReadLock(func(s2 SetLike[T]) {
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
			la.WithReadLock(func(a2 SetLike[T]) {
				lb.WithReadLock(func(b2 SetLike[T]) {
					r = fn(a2, b2)
				})
			})
			return r
		}
	}
	return fn(a, b)
}
