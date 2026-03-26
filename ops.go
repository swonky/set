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

// AccumulateTry returns a sequence of intermediate accumulation results.
// Starting from the first set, fn is applied pairwise across sets.
// Iteration stops early if fn returns false or the consumer stops iteration.
// For zero sets, a single empty set is yielded.
// For one set, a clone of that set is yielded.
func AccumulateTry[T comparable](
	sets []Set[T],
	fn func(a, b Set[T]) (Set[T], bool),
) iter.Seq[Set[T]] {
	return func(yield func(Set[T]) bool) {
		if len(sets) == 0 {
			yield(New[T]())
			return
		}

		r := sets[0].Clone()

		for _, s := range sets[1:] {
			var ok bool
			r, ok = fn(r, s)
			if !yield(r) || !ok {
				return
			}
		}
	}
}

// Accumulate is like AccumulateTry but always continues until all sets are processed.
func Accumulate[T comparable](sets []Set[T], fn func(a, b Set[T]) Set[T]) iter.Seq[Set[T]] {
	return AccumulateTry(sets, func(a, b Set[T]) (Set[T], bool) { return fn(a, b), true })
}

// ReduceTry reduces sets using fn, which may stop early by returning false.
// The last successfully produced value is returned.
// If no sets are provided, an empty set is returned.
func ReduceTry[T comparable](
	sets []Set[T],
	fn func(a, b Set[T]) (Set[T], bool),
) Set[T] {
	if len(sets) == 0 {
		return New[T]()
	}
	var last Set[T]
	for v := range AccumulateTry(sets, fn) {
		last = v
	}
	return last
}

// ReduceWhile reduces sets using fn and stops when pred returns false.
// The last value satisfying the predicate is returned.
// If no sets are provided, an empty set is returned.
func ReduceWhile[T comparable](
	sets []Set[T],
	fn func(a, b Set[T]) Set[T],
	pred func(s Set[T]) bool,
) Set[T] {
	if len(sets) == 0 {
		return New[T]()
	}
	var last Set[T]
	for v := range Accumulate(sets, fn) {
		if last = v; !pred(last) {
			break
		}
	}
	return last
}

// Reduce reduces sets using fn until all sets are processed.
// The final accumulated value is returned.
// If no sets are provided, an empty set is returned.
func Reduce[T comparable](
	sets []Set[T],
	fn func(a, b Set[T]) Set[T],
) Set[T] {
	if len(sets) == 0 {
		return New[T]()
	}
	var last Set[T]
	for v := range Accumulate(sets, fn) {
		last = v
	}
	return last
}

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

// Clone returns a shallow copy of the set.
func Clone[T comparable](s Set[T]) Set[T] {
	return s.Clone()
}

// Chain returns a sequence of unique elements from the provided sets,
// preserving first occurrence order across sets.
func Chain[T comparable](sets ...Set[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := New[T]()
		for _, s := range sets {
			for e := range s.Iter() {
				if ok := seen.Add(e); !ok {
					if !yield(e) {
						return
					}
				}
			}
		}
	}
}
