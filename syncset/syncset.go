package syncset

import (
	"iter"
	"maps"
	"sync"

	"github.com/swonky/set"
	"github.com/swonky/set/internal/base"
)

var (
	_ set.SetLike[int]     = (*SyncSet[int])(nil)
	_ set.MutableSet[int]  = (*SyncSet[int])(nil)
	_ set.LockableSet[int] = (*SyncSet[int])(nil)
	_ set.AsSetter[int]    = (*SyncSet[int])(nil)
)

// SyncSet is a concurrency-safe wrapper around Set.
// All operations are guarded by a RWMutex.
// Iteration holds a read lock for the duration of the sequence.
type SyncSet[T comparable] struct {
	mu     sync.RWMutex
	values set.Set[T]
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
func (s *SyncSet[T]) Clone() *SyncSet[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &SyncSet[T]{values: maps.Clone(s.values)}
}

// Len returns the number of elements in the set.
func (s *SyncSet[T]) Len() int {
	s.mu.RLock()
	n := len(s.values)
	s.mu.RUnlock()
	return n
}

// Contains reports whether item exists in the set.
func (s *SyncSet[T]) Contains(item T) bool {
	s.mu.RLock()
	_, ok := s.values[item]
	s.mu.RUnlock()
	return ok
}

// Add inserts item into the set.
func (s *SyncSet[T]) Add(item T) {
	s.mu.Lock()
	s.values[item] = struct{}{}
	s.mu.Unlock()
}

// Delete removes item from the set.
func (s *SyncSet[T]) Delete(item T) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	delete(s.values, item)
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
func (s *SyncSet[T]) WithRLock(fn func(set.SetLike[T])) {
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
func (s *SyncSet[T]) WithLock(fn func(set.MutableSet[T])) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(s.values)
}

func (s *SyncSet[T]) AsSet() set.Set[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return maps.Clone(s.values)
}
