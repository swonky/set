package set

import "iter"

// New returns a new Set containing the provided items.
// If no items are provided, it returns an empty set.
func New[T comparable](items ...T) Set[T] {
	s := make(Set[T], len(items))
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

// Collect returns a new Set containing all elements produced by the iterator.
// If a positive capacity hint is provided, it is used to preallocate the set.
func Collect[T comparable](it iter.Seq[T], capacity ...int) Set[T] {
	var s Set[T]
	if len(capacity) > 0 && capacity[0] > 0 {
		s = make(Set[T], capacity[0])
	} else {
		s = make(Set[T])
	}
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
func UnionIter[T comparable](sets ...Set[T]) iter.Seq[T] {
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
func IntersectIter[T comparable](sets ...Set[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		if len(sets) == 0 {
			return
		}

		// Find smallest set index
		smallestIdx := 0
		for i := 1; i < len(sets); i++ {
			if len(sets[i]) < len(sets[smallestIdx]) {
				smallestIdx = i
			}
		}

		smallest := sets[smallestIdx]

		for k := range smallest {
			ok := true
			for i, s := range sets {
				if i == smallestIdx {
					continue
				}
				if _, exists := s[k]; !exists {
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
