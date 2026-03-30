package set

import (
	"iter"
	"slices"
)

type SetLike[T comparable] interface {
	All(fn func(T) bool) bool
	Any(fn func(T) bool) bool
	AsSlice() []T
	Find(fn func(T) bool) (T, bool)
	First() (T, bool)
	Has(item T) bool
	HasAll(item ...T) bool
	HasAny(item ...T) bool
	IsEmpty() bool
	Range(func(T) bool)
	Len() int
	String() string
}

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

// Clone returns a shallow copy of the set.
func Clone[T comparable](s Set[T]) Set[T] {
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

func IntersectIter[T comparable](sets []SetLike[T]) iter.Seq[T] {
	return func(yield func(T) bool) { IntersectRange(sets, yield) }
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
func Transform[T, C comparable](s SetLike[T], pred func(T) C) map[C]Set[T] {
	if pred == nil {
		panic("nil predicate in set.Transform")
	}
	r := make(map[C]Set[T])
	for k := range s.Range {
		c := pred(k)
		g := r[c]
		if g == nil {
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
// Use IntersectAll when you need to retain the full result.
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
func IntersectRange[T comparable](sets []SetLike[T], fn func(T) bool) {
	switch len(sets) {
	case 0:
		return
	case 1:
		sets[0].Range(fn)
		return
	}

	smallest := 0
	for i := 1; i < len(sets); i++ {
		if sets[i].Len() < sets[smallest].Len() {
			smallest = i
		}
	}

	s0 := sets[smallest]

	switch s := any(s0).(type) {
	case Set[T]:
		for k := range s {
			if !intersectHasAll(sets, smallest, k) {
				continue
			}
			if !fn(k) {
				return
			}
		}
		return

	case FrozenSet[T]:
		for k := range s.s {
			if !intersectHasAll(sets, smallest, k) {
				continue
			}
			if !fn(k) {
				return
			}
		}
		return

	case SyncSet[T]:
		s.mu.RLock()
		for k := range s.s {
			if !intersectHasAll(sets, smallest, k) {
				continue
			}
			if !fn(k) {
				s.mu.RUnlock()
				return
			}
		}
		s.mu.RUnlock()
		return
	}

	s0.Range(func(k T) bool {
		if !intersectHasAll(sets, smallest, k) {
			return true
		}
		return fn(k)
	})
}

func intersectHasAll[T comparable](sets []SetLike[T], skip int, k T) bool {
	for i := range sets {
		if i == skip {
			continue
		}
		if !sets[i].Has(k) {
			return false
		}
	}
	return true
}
