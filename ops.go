package set

import "iter"

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
	Iter() iter.Seq[T]
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
			for e := range s.Iter() {
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

// IntersectIter returns a lazy sequence of elements present in all provided sets.
// It iterates the smallest input set and yields elements that exist in every other set.
// Iteration stops early if the consumer returns false.
//
// No intermediate sets are allocated. Membership checks are performed on demand.
// For zero sets, the sequence yields nothing.
func IntersectIter[T comparable](sets ...SetLike[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		if len(sets) == 0 {
			return
		}

		// Find smallest set index
		smallestIdx := 0
		for i := 1; i < len(sets); i++ {
			if sets[i].Len() < sets[smallestIdx].Len() {
				smallestIdx = i
			}
		}

		smallest := sets[smallestIdx]

		for k := range smallest.Iter() {
			ok := true
			for i, s := range sets {
				if i == smallestIdx {
					continue
				}
				if !s.Has(k) {
					ok = false
					break
				}
			}
			if ok {
				if !yield(k) {
					return
				}
			}
		}
	}
}
