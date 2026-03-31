// Package set provides a generic set type and operations for working with sets of comparable elements.
// A set is implemented as a map[T]struct{} for memory efficiency and O(1) average-case operations.
package set

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
)

var _ SetLike[int] = Set[int]{}

// Set[T] represents a mathematical set of comparable elements.
// The zero value is ready to use but prefer using New() for initialization.
type Set[T comparable] map[T]struct{}

func (s Set[T]) Freeze() FrozenSet[T] {
	return FrozenSet[T]{s: maps.Clone(s)}
}

// Union returns a new set containing all elements from s and o.
func (s Set[T]) Union(o Set[T]) Set[T] {
	if len(s) > len(o) {
		s, o = o, s
	}
	r := s.Clone()
	for k := range o {
		r[k] = struct{}{}
	}
	return r
}

// UnionIter returns a lazy sequence of elements from the larger set,
// followed by elements from the smaller set not already present.
// Membership checks are performed against the smaller set to reduce probes.
// No intermediate state is allocated.
func (s Set[T]) UnionIter(o Set[T]) iter.Seq[T] {
	large, small := s, o
	if len(large) < len(small) {
		large, small = small, large
	}

	return func(yield func(T) bool) {
		for k := range large {
			if !yield(k) {
				return
			}
		}
		for k := range small {
			if _, exists := large[k]; !exists {
				if !yield(k) {
					return
				}
			}
		}
	}
}

// Intersect returns a new [Set] containing elements present in both s and o.
// Neither s or o are mutated.
func (s Set[T]) Intersect(o Set[T]) Set[T] {
	if len(s) > len(o) {
		s, o = o, s
	}
	out := make(Set[T], len(s))
	for k := range s {
		if _, ok := o[k]; ok {
			out[k] = struct{}{}
		}
	}
	return out
}

// IntersectIter returns a lazy sequence of elements present in both s and o.
// The smaller set is selected before iteration to avoid branching in the hot path.
// No intermediate set is allocated and iteration may stop early.
func (s Set[T]) IntersectIter(o Set[T]) iter.Seq[T] {
	small, large := s, o
	if len(small) > len(large) {
		small, large = large, small
	}

	return func(yield func(T) bool) {
		for k := range small {
			if _, exists := large[k]; exists {
				if !yield(k) {
					return
				}
			}
		}
	}
}

// Diff returns elements in s but not in o.
func (s Set[T]) Diff(o Set[T]) Set[T] {
	out := make(Set[T], len(s))
	for k := range s {
		if _, ok := o[k]; !ok {
			out[k] = struct{}{}
		}
	}
	return out
}

// SymDiff returns elements in either set but not both.
func (s Set[T]) SymDiff(o Set[T]) Set[T] {
	return s.Diff(o).Union(o.Diff(s))
}

// Has reports whether the item is present.
func (s Set[T]) Has(item T) bool {
	_, ok := s[item]
	return ok
}

// HasAny reports whether any of the item are present.
func (s Set[T]) HasAny(item ...T) bool {
	return slices.ContainsFunc(item, s.Has)
}

// HasAll reports whether all of the item are present.
func (s Set[T]) HasAll(item ...T) bool {
	for _, v := range item {
		if !s.Has(v) {
			return false
		}
	}
	return true
}

// Clone returns a shallow copy.
func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}

// Len returns the number of elements.
func (s Set[T]) Len() int {
	return len(s)
}

// IsEmpty reports whether the set is empty.
func (s Set[T]) IsEmpty() bool {
	return len(s) == 0
}

// IsSubsetOf reports whether s ⊆ o.
func (s Set[T]) IsSubsetOf(o Set[T]) bool {
	for k := range s {
		if !o.Has(k) {
			return false
		}
	}
	return true
}

// IsSupersetOf reports whether s ⊇ o.
func (s Set[T]) IsSupersetOf(o Set[T]) bool {
	return o.IsSubsetOf(s)
}

// Equal reports whether two sets contain the same elements.
func (s Set[T]) Equal(o Set[T]) bool {
	return len(s) == len(o) && s.IsSubsetOf(o)
}

// Range calls yield for each element in the set.
//
// Iteration continues until all elements have been processed or yield returns false.
//
// Iteration order is not specified and may vary between calls.
//
// Range performs no allocations.
//
// The set must not be modified during iteration.
//
// A nil set produces no calls to yield.
//
// Range panics if yield is nil.
func (s Set[T]) doRange(yield func(T) bool) {
	if yield == nil {
		panic("nil yield function in Set[T].Range")
	}
	for k := range s {
		if !yield(k) {
			return
		}
	}
}

func (s Set[T]) Range(yield func(T) bool) {
	s.doRange(yield)
}

// AsSlice returns elements as a slice.
func (s Set[T]) AsSlice() []T {
	out := make([]T, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	return out
}

// String returns a string representation.
func (s Set[T]) String() string {
	elems := make([]string, 0, len(s))
	for k := range s {
		elems = append(elems, fmt.Sprintf("%v", k))
	}
	return "{" + strings.Join(elems, ", ") + "}"
}

// Filter returns a new set containing elements satisfying fn.
func (s Set[T]) Filter(fn func(T) bool) Set[T] {
	out := make(Set[T], len(s))
	for k := range s {
		if fn(k) {
			out[k] = struct{}{}
		}
	}
	return out
}

// Any reports whether any element satisfies fn.
func (s Set[T]) Any(fn func(T) bool) bool {
	for k := range s {
		if fn(k) {
			return true
		}
	}
	return false
}

// All reports whether all elements satisfy fn.
func (s Set[T]) All(fn func(T) bool) bool {
	for k := range s {
		if !fn(k) {
			return false
		}
	}
	return true
}

// Find returns the first element satisfying fn.
func (s Set[T]) Find(fn func(T) bool) (T, bool) {
	for k := range s {
		if fn(k) {
			return k, true
		}
	}
	var zero T
	return zero, false
}

// First returns an arbitrary element.
func (s Set[T]) First() (T, bool) {
	for k := range s {
		return k, true
	}
	var zero T
	return zero, false
}

// Mutable operations

// Partition
func (s Set[T]) Partition(pred func(T) bool) (Set[T], Set[T]) {
	if pred == nil {
		panic("nil predicate")
	}
	a := make(Set[T], len(s)/2)
	b := make(Set[T], len(s)/2)
	for k := range s {
		if pred(k) {
			a[k] = struct{}{}
		} else {
			b[k] = struct{}{}
		}
	}
	return a, b
}

// Mutable operations

// UnionInto inserts all elements from o into s in place.
// It performs a single pass over o with no allocations.
func (s Set[T]) UnionInto(o Set[T]) {
	for k := range o {
		s[k] = struct{}{}
	}
}

// AddAll inserts one or more items into the set.
func (s Set[T]) AddAll(items ...T) {
	for _, item := range items {
		s[item] = struct{}{}
	}
}

// Add inserts an item into the set.
func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

// AddCheck inserts an item into the set. Returns true if item was already present.
func (s Set[T]) AddCheck(item T) bool {
	if _, ok := s[item]; ok {
		return true
	}
	s[item] = struct{}{}
	return false
}

// Clear removes all elements.
func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Delete removes an item from the set.
func (s Set[T]) Delete(item T) {
	delete(s, item)
}

// Pop removes and returns an arbitrary element.
func (s Set[T]) Pop() (T, bool) {
	for k := range s {
		delete(s, k)
		return k, true
	}
	var zero T
	return zero, false
}
