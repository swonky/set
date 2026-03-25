package set

import "iter"

// New creates a new Set containing the provided items.
// If no items are provided, returns an empty set.
func New[T comparable](items ...T) Set[T] {
	s := make(Set[T], len(items))
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

// FromIter creates a new Set from an iterator sequence.
// An optional capacity hint may be provided to reduce reallocations.
func FromIter[T comparable](it iter.Seq[T], capacity ...int) Set[T] {
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

// UnionAll returns a new Set containing all elements from the provided sets.
func UnionAll[T comparable](sets ...Set[T]) Set[T] {
	switch len(sets) {
	case 0:
		return New[T]()
	case 1:
		return sets[0].Clone()
	}
	r := sets[0].Clone()
	for _, s := range sets[1:] {
		r = r.Union(s)
	}
	return r
}

// IntersectAll returns a new Set containing elements present in ALL input sets.
func IntersectAll[T comparable](sets ...Set[T]) Set[T] {
	if len(sets) == 0 {
		return New[T]()
	}
	r := sets[0].Clone()
	for _, s := range sets[1:] {
		r = r.Intersect(s)
		if r.IsEmpty() {
			break
		}
	}
	return r
}
