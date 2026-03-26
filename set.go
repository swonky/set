// Package set provides a generic Set type and operations for working with sets of comparable elements.
// A Set is implemented as a map[T]struct{} for memory efficiency and O(1) average-case operations.
package set

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
)

// Set represents a mathematical set of comparable elements.
// The zero value is ready to use but prefer using New() for initialization.
type Set[T comparable] map[T]struct{}

// Union returns a new Set containing all elements from both sets.
func (s Set[T]) Union(o Set[T]) Set[T] {
	if len(o) > len(s) {
		s, o = o, s
	}
	r := s.Clone()
	for k := range o {
		r[k] = struct{}{}
	}
	return r
}

// UnionInto inserts all elements from the other set into this set.
func (s Set[T]) UnionInto(o Set[T]) {
	for k := range o {
		s[k] = struct{}{}
	}
}

// Intersect returns a new Set containing elements present in both sets.
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

// SymmetricDiff returns elements in either set but not both.
func (s Set[T]) SymmetricDiff(o Set[T]) Set[T] {
	return s.Diff(o).Union(o.Diff(s))
}

// AddAll inserts one or more items into the set.
func (s Set[T]) AddAll(items ...T) {
	for _, item := range items {
		s[item] = struct{}{}
	}
}

// Add inserts an item into the set. Returns true if item was already present.
func (s Set[T]) Add(item T) bool {
	if s.Has(item) {
		return true
	}
	s[item] = struct{}{}
	return false
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

// Clear removes all elements.
func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
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

// Iter returns an iterator over the set.
func (s Set[T]) Iter() iter.Seq[T] {
	return maps.Keys(s)
}

// AsSlice returns elements as a slice.
func (s Set[T]) AsSlice() []T {
	return slices.Collect(s.Iter())
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
