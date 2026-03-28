// Package set provides a generic `Set` type for working with unordered collections of unique, comparable values.
//
// A Set is implemented as a map[T]struct{}, offering O(1) average-case performance for insertion, deletion, and membership checks.
//
// Basic usage:
//
//	s := set.New(1, 2, 3)
//	s.Add(4)
//
//	if s.Has(2) {
//		fmt.Println("contains 2")
//	}
//
// Set operations:
//
//	a := set.New(1, 2, 3)
//	b := set.New(3, 4, 5)
//
//	union := a.Union(b)        // {1, 2, 3, 4, 5}
//	inter := a.Intersect(b)    // {3}
//	diff := a.Diff(b)          // {1, 2}
//
// Mutating operations:
//
//	a.UnionInto(b) // a now contains all elements from both sets
//
// Iteration:
//
//	for v := range s.Iter() {
//		fmt.Println(v)
//	}
//
// Predicates:
//
//	s.Any(func(v int) bool { return v%2 == 0 })
//	s.All(func(v int) bool { return v > 0 })
//
// Construction from iterators:
//
//	s := set.Collect(iter.Seq[int]{...})
//
// Reduction and accumulation:
//
//	sets := []set.Set[int]{
//		set.New(1, 2),
//		set.New(2, 3),
//		set.New(3, 4),
//	}
//
// Reduce combines sets left-to-right using a binary operation.
//
//	r := set.Reduce(sets, func(a, b set.Set[int]) set.Set[int] {
//		return a.Union(b)
//	}) // {1,2,3,4}
//
// ReduceWhile stops when the predicate fails.
//
//	r = set.ReduceWhile(
//		sets,
//		func(a, b set.Set[int]) set.Set[int] { return a.Union(b) },
//		func(s set.Set[int]) bool { return s.Len() < 4 },
//	)
//
// ReduceTry allows the reducer to control early termination.
//
//	r = set.ReduceTry(
//		sets,
//		func(a, b set.Set[int]) (set.Set[int], bool) {
//			r := a.Intersect(b)
//			return r, !r.IsEmpty()
//		},
//	)
//
// Accumulate produces intermediate results.
//
//	for v := range set.Accumulate(sets, func(a, b set.Set[int]) set.Set[int] {
//		return a.Union(b)
//	}) {
//		fmt.Println(v)
//	}
//
// The iteration order of a Set is not defined and may vary between runs.
//
// This package is designed to be small, predictable, and idiomatic, providing a
// practical set abstraction without additional dependencies or complex semantics.
package set
