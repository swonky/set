package syncset

import (
	"iter"
	"maps"
	"sync"
	"unsafe"

	"github.com/swonky/set/internal/base"
)

var _ base.SetLike[int] = (*SyncSet[int])(nil)

// SyncSet is a concurrency-safe wrapper around Set.
// All operations are guarded by a RWMutex.
// Iteration holds a read lock for the duration of the sequence.
type SyncSet[T comparable] struct {
	mu   sync.RWMutex
	smap base.Set[T]
}

// NewSync returns an empty SyncSet.
func New[T comparable](cap ...int) *SyncSet[T] {
	return &SyncSet[T]{
		smap: make(map[T]struct{}, base.GetCap(cap...)),
	}
}

// From
func From[T comparable](elems ...T) *SyncSet[T] {
	s := New[T](len(elems))
	for _, t := range elems {
		s.smap[t] = struct{}{}
	}
	return s
}

// Collect
func Collect[T comparable](seq iter.Seq[T]) *SyncSet[T] {
	s := New[T]()
	for t := range seq {
		s.smap[t] = struct{}{}
	}
	return s
}

func (ss *SyncSet[T]) Clone() *SyncSet[T] {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return &SyncSet[T]{smap: maps.Clone(ss.smap)}
}

// Len returns the number of elements in the set.
func (ss *SyncSet[T]) Len() int {
	ss.mu.RLock()
	n := len(ss.smap)
	ss.mu.RUnlock()
	return n
}

// Contains reports whether item exists in the set.
func (ss *SyncSet[T]) Contains(item T) bool {
	ss.mu.RLock()
	_, ok := ss.smap[item]
	ss.mu.RUnlock()
	return ok
}

// Add inserts item into the set.
func (ss *SyncSet[T]) Add(item T) {
	ss.mu.Lock()
	ss.smap[item] = struct{}{}
	ss.mu.Unlock()
}

// Delete removes item from the set.
func (ss *SyncSet[T]) Delete(item T) {
	ss.mu.Lock()
	delete(ss.smap, item)
	ss.mu.Unlock()
}

// Iter returns a sequence of elements in the set.
// A read lock is held for the duration of iteration.
func (ss *SyncSet[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		ss.mu.RLock()
		defer ss.mu.RUnlock()
		for k := range ss.smap {
			if !yield(k) {
				return
			}
		}
	}
}

func (s *SyncSet[T]) Range(fn func(T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k := range s.smap {
		if !fn(k) {
			return
		}
	}
}

func (s *SyncSet[T]) WithReadLock(fn func(base.SetLike[T])) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn(s.smap)
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
