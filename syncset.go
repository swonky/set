package set

import (
	"iter"
	"sync"
	"unsafe"
)

var _ SetLike[int] = (*SyncSet[int])(nil)

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

func (ss *SyncSet[T]) Clone() FrozenSet[T] {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.s.Clone().Freeze()
}

func (ss *SyncSet[T]) Union(o *SyncSet[T]) FrozenSet[T] {
	lock2(ss, o)
	r := ss.s.Union(o.s)
	unlock2(ss, o)
	return r.Freeze()
}

func (ss *SyncSet[T]) Intersect(o *SyncSet[T]) FrozenSet[T] {
	lock2(ss, o)
	r := ss.s.Intersect(o.s)
	unlock2(ss, o)
	return r.Freeze()
}

func (ss *SyncSet[T]) Diff(o *SyncSet[T]) FrozenSet[T] {
	lock2(ss, o)
	r := ss.s.Diff(o.s)
	unlock2(ss, o)
	return r.Freeze()
}

func (ss *SyncSet[T]) SymDiff(o *SyncSet[T]) FrozenSet[T] {
	lock2(ss, o)
	r := ss.s.SymDiff(o.s)
	unlock2(ss, o)
	return r.Freeze()
}

func (ss *SyncSet[T]) Filter(fn func(T) bool) FrozenSet[T] {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.s.Filter(fn).Freeze()
}

func (ss *SyncSet[T]) Equal(o *SyncSet[T]) bool {
	lock2(ss, o)
	r := ss.s.Equal(o.s)
	unlock2(ss, o)
	return r
}

func (ss *SyncSet[T]) IsSubsetOf(o *SyncSet[T]) bool {
	lock2(ss, o)
	r := ss.s.IsSubsetOf(o.s)
	unlock2(ss, o)
	return r
}

func (ss *SyncSet[T]) IsSupersetOf(o *SyncSet[T]) bool {
	lock2(ss, o)
	r := ss.s.IsSupersetOf(o.s)
	unlock2(ss, o)
	return r
}

// AsSet creates a new mutable thread-unsafe Set.
// The new Set is initialized with a clone of the contents of ss.
func (ss *SyncSet[T]) AsSet() Set[T] {
	ss.mu.RLock()
	s := ss.s.Clone()
	ss.mu.RUnlock()
	return s
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

func (ss *SyncSet[T]) UnionInto(o *SyncSet[T]) {
	if ss == o {
		return
	}

	if uintptr(unsafe.Pointer(ss)) < uintptr(unsafe.Pointer(o)) {
		ss.mu.Lock()
		o.mu.RLock()
	} else {
		o.mu.RLock()
		ss.mu.Lock()
	}

	for k := range o.s {
		ss.s[k] = struct{}{}
	}

	ss.mu.Unlock()
	o.mu.RUnlock()
}

// UnionIter returns an iterator of all elements from both sets.
// Locks are held only during snapshot creation.
func (ss *SyncSet[T]) UnionIter(o *SyncSet[T]) iter.Seq[T] {
	return ss.Clone().UnionIter(o.Clone())
}

// IntersectIter returns an iterator of elements present in both sets.
// Locks are held only during snapshot creation.
func (ss *SyncSet[T]) IntersectIter(o *SyncSet[T]) iter.Seq[T] {
	a := ss.Clone()
	b := o.Clone()
	return a.IntersectIter(b)
}

func (ss *SyncSet[T]) All(fn func(T) bool) bool {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.s.All(fn)
}

func (ss *SyncSet[T]) Any(fn func(T) bool) bool {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.s.Any(fn)
}

func (ss *SyncSet[T]) AsSlice() []T {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.s.AsSlice()
}

func (ss *SyncSet[T]) Find(fn func(T) bool) (T, bool) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.s.Find(fn)
}

func (ss *SyncSet[T]) First() (T, bool) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.s.First()
}

func (ss *SyncSet[T]) HasAll(item ...T) bool {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.s.HasAll(item...)
}

func (ss *SyncSet[T]) HasAny(item ...T) bool {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.s.HasAny(item...)
}

func (ss *SyncSet[T]) IsEmpty() bool {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.s.IsEmpty()
}

func (ss *SyncSet[T]) String() string {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.s.String()
}

func (s *SyncSet[T]) Range(fn func(T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k := range s.s {
		if !fn(k) {
			return
		}
	}
}

func lock2[T comparable](a, b *SyncSet[T]) {
	if a == b {
		a.mu.RLock()
		return
	}
	if uintptr(unsafe.Pointer(a)) < uintptr(unsafe.Pointer(b)) {
		a.mu.RLock()
		b.mu.RLock()
	} else {
		b.mu.RLock()
		a.mu.RLock()
	}
}

func unlock2[T comparable](a, b *SyncSet[T]) {
	if a == b {
		a.mu.RUnlock()
		return
	}
	a.mu.RUnlock()
	b.mu.RUnlock()
}
