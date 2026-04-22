package syncset

import (
	"iter"
	"maps"
	"sync"

	"github.com/swonky/set/internal/base"
)

var (
	_ base.SetLike[int]     = (*SyncSet[int])(nil)
	_ base.MutableSet[int]  = (*SyncSet[int])(nil)
	_ base.LockableSet[int] = (*SyncSet[int])(nil)
	_ base.AsSetter[int]    = (*SyncSet[int])(nil)
)

// SyncSet is a concurrency-safe wrapper around Set.
// All operations are guarded by a RWMutex.
// Iteration holds a read lock for the duration of the sequence.
type SyncSet[T comparable] struct {
	mu     sync.RWMutex
	values base.Set[T]
}

// New
func New[T comparable](cap ...int) *SyncSet[T] {
	return &SyncSet[T]{
		values: make(map[T]struct{}, base.GetCap(cap...)),
	}
}

// From
func From[T comparable](elems ...T) *SyncSet[T] {
	s := New[T](len(elems))
	for _, t := range elems {
		s.values[t] = struct{}{}
	}
	return s
}

// Collect
func Collect[T comparable](seq iter.Seq[T]) *SyncSet[T] {
	s := New[T]()
	for t := range seq {
		s.values[t] = struct{}{}
	}
	return s
}

// Clone
func (ss *SyncSet[T]) Clone() *SyncSet[T] {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return &SyncSet[T]{values: maps.Clone(ss.values)}
}

// Len returns the number of elements in the set.
func (ss *SyncSet[T]) Len() int {
	ss.mu.RLock()
	n := len(ss.values)
	ss.mu.RUnlock()
	return n
}

// Contains reports whether item exists in the set.
func (ss *SyncSet[T]) Contains(item T) bool {
	ss.mu.RLock()
	_, ok := ss.values[item]
	ss.mu.RUnlock()
	return ok
}

// Add inserts item into the set.
func (ss *SyncSet[T]) Add(item T) {
	ss.mu.Lock()
	ss.values[item] = struct{}{}
	ss.mu.Unlock()
}

// Delete removes item from the set.
func (ss *SyncSet[T]) Delete(item T) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	delete(ss.values, item)
}

// Range implements iter.Seq.
// It calls yield for each element in the set while holding a read lock.
// Iteration stops when fn returns false.
// The iteration order is unspecified.
// fn must not call methods on s.
//
// Use [SyncSet.WithRLock] for non-iterative read operations or when multiple reads
// should be performed while holding a single lock.
func (s *SyncSet[T]) Range(yield func(T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for k := range s.values {
		if !yield(k) {
			return
		}
	}
}

// WithRLock calls fn while holding a read lock and passes a temporary
// read-only view of the set for use only during the callback.
//
// The view must not be retained after fn returns.
// fn must not call methods on s.
//
// Use [SyncSet.WithRWLock] when write operations are required.
// Use [SyncSet.Range] when only element iteration is required.
func (s *SyncSet[T]) WithRLock(fn func(base.SetLike[T])) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn(s.values)
}

// WithLock calls fn while holding a read/write lock and passes a temporary
// view of the set for use only during the callback.
//
// The view must not be retained after fn returns.
// fn must not call methods on s.
//
// Use [SyncSet.WithRLock] when only read operations are required.
// Use [SyncSet.Range] when only element iteration is required.
func (s *SyncSet[T]) WithLock(fn func(base.Set[T])) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(s.values)
}

func (ss *SyncSet[T]) AsSet() base.Set[T] {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	return maps.Clone(ss.values)
}
