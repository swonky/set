package set

import (
	"iter"
	"sync"
)

// SyncSet is a concurrency-safe wrapper around Set.
// All operations are guarded by a RWMutex.
// Iteration holds a read lock for the duration of the sequence.
type SyncSet[T comparable] struct {
	mu sync.RWMutex
	s  Set[T]
}

// NewSync returns an empty SyncSet.
func NewSync[T comparable]() *SyncSet[T] {
	return &SyncSet[T]{s: New[T]()}
}

// FromSet creates a SyncSet initialized with a clone of s.
func FromSet[T comparable](s Set[T]) *SyncSet[T] {
	return &SyncSet[T]{s: s.Clone()}
}

// Len returns the number of elements in the set.
func (ss *SyncSet[T]) Len() int {
	ss.mu.RLock()
	n := len(ss.s)
	ss.mu.RUnlock()
	return n
}

// Has reports whether item exists in the set.
func (ss *SyncSet[T]) Has(item T) bool {
	ss.mu.RLock()
	_, ok := ss.s[item]
	ss.mu.RUnlock()
	return ok
}

// Add inserts item into the set.
func (ss *SyncSet[T]) Add(item T) {
	ss.mu.Lock()
	ss.s[item] = struct{}{}
	ss.mu.Unlock()
}

// HasAdd inserts item and reports whether it was already present.
func (ss *SyncSet[T]) HasAdd(item T) bool {
	ss.mu.Lock()
	_, existed := ss.s[item]
	ss.s[item] = struct{}{}
	ss.mu.Unlock()
	return existed
}

// Delete removes item from the set.
func (ss *SyncSet[T]) Delete(item T) {
	ss.mu.Lock()
	delete(ss.s, item)
	ss.mu.Unlock()
}

// Clear removes all elements from the set.
func (ss *SyncSet[T]) Clear() {
	ss.mu.Lock()
	ss.s = New[T]()
	ss.mu.Unlock()
}

// Clone returns a snapshot copy of the underlying Set.
func (ss *SyncSet[T]) Clone() Set[T] {
	ss.mu.RLock()
	c := ss.s.Clone()
	ss.mu.RUnlock()
	return c
}

// Iter returns a sequence of elements in the set.
// A read lock is held for the duration of iteration.
func (ss *SyncSet[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		ss.mu.RLock()
		defer ss.mu.RUnlock()
		for k := range ss.s {
			if !yield(k) {
				return
			}
		}
	}
}

// Union returns a new Set containing all elements from both sets.
// Read locks are held only for the duration of snapshotting.
func (ss *SyncSet[T]) Union(o *SyncSet[T]) Set[T] {
	a := ss.Clone()
	b := o.Clone()
	return a.Union(b)
}

// UnionInto inserts all elements from o into ss.
func (ss *SyncSet[T]) UnionInto(o *SyncSet[T]) {
	o.mu.RLock()
	ss.mu.Lock()
	for k := range o.s {
		ss.s[k] = struct{}{}
	}
	ss.mu.Unlock()
	o.mu.RUnlock()
}

// UnionIter returns an iterator of all elements from both sets.
// Locks are held only during snapshot creation.
func (ss *SyncSet[T]) UnionIter(o *SyncSet[T]) iter.Seq[T] {
	a := ss.Clone()
	b := o.Clone()
	return a.UnionIter(b)
}

// Intersect returns a new Set containing elements present in both sets.
// Read locks are held only for the duration of snapshotting.
func (ss *SyncSet[T]) Intersect(o *SyncSet[T]) *SyncSet[T] {
	return FromSet(ss.Clone().Intersect(o.Clone()))
}

// IntersectIter returns an iterator of elements present in both sets.
// Locks are held only during snapshot creation.
func (ss *SyncSet[T]) IntersectIter(o *SyncSet[T]) iter.Seq[T] {
	a := ss.Clone()
	b := o.Clone()
	return a.IntersectIter(b)
}
