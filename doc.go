// Package set provides a generic Set type for working with unordered collections
// of unique, comparable values.
//
// A Set is implemented as a map[T]struct{}, offering O(1) average-case performance
// for insertion, deletion, and membership checks. The API follows a consistent model:
//
//   - Pure operations return new sets (e.g. Union, Intersect, Diff)
//   - Mutating operations modify the receiver in place (e.g. Add, Delete, UnionInto)
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
//	s := set.FromIter(iter.Seq[int]{...})
//
// The iteration order of a Set is not defined and may vary between runs.
//
// This package is designed to be small, predictable, and idiomatic, providing a
// practical set abstraction without additional dependencies or complex semantics.
package set
