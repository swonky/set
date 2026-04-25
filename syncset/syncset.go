package syncset

import (
	"iter"
	"sync"

	"github.com/swonky/set"
	"github.com/swonky/set/internal/base"
)

var (
	_ set.SetLike[int]                   = (*SyncSet[set.Set[int], int])(nil)
	_ set.MutableSet[int]                = (*SyncSet[set.Set[int], int])(nil)
	_ set.LockableSet[set.Set[int], int] = (*SyncSet[set.Set[int], int])(nil)
)

// SyncSet is a concurrency-safe wrapper around Set.
// All operations are guarded by a RWMutex.
// Iteration holds a read lock for the duration of the sequence.
type SyncSet[S set.SetLike[T], T comparable] struct {
	mu     sync.RWMutex
	values S
}

// New
func New[T comparable](cap ...int) *SyncSet[set.Set[T], T] {
	return &SyncSet[set.Set[T], T]{
		values: make(map[T]struct{}, base.GetCap(cap...)),
	}
}

// From
func From[T comparable](elems ...T) *SyncSet[set.Set[T], T] {
	s := New[T](len(elems))
	for _, t := range elems {
		s.values[t] = struct{}{}
	}
	return s
}

// Collect
func Collect[T comparable](seq iter.Seq[T]) *SyncSet[set.Set[T], T] {
	s := New[T]()
	for t := range seq {
		s.values[t] = struct{}{}
	}
	return s
}

// Len returns the number of elements in the set.
func (s *SyncSet[S, T]) Len() int {
	s.mu.RLock()
	n := s.values.Len()
	s.mu.RUnlock()
	return n
}

// Contains reports whether item exists in the set.
func (s *SyncSet[S, T]) Contains(item T) bool {
	s.mu.RLock()
	ok := s.values.Contains(item)
	s.mu.RUnlock()
	return ok
}

// Add inserts item into the set.
func (s *SyncSet[S, T]) Add(item T) {
	s.mu.Lock()
	s.values.Add(item)
	s.mu.Unlock()
}

// Delete removes item from the set.
func (s *SyncSet[S, T]) Delete(item T) {
	s.mu.RLock()
	s.values.Delete(item)
	s.mu.RUnlock()
}

// Range implements iter.Seq.
// It calls yield for each element in the set while holding a read lock.
// Iteration stops when fn returns false.
// The iteration order is unspecified.
// fn must not call methods on s.
//
// Use [SyncSet.WithRLock] for non-iterative read operations or when multiple reads
// should be performed while holding a single lock.
func (s *SyncSet[S, T]) Range(yield func(T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.values.Range(func(t T) bool { return yield(t) })
}

// WithRLock calls fn while holding a read lock and passes a temporary
// read-only view of the set for use only during the callback.
//
// The view must not be retained after fn returns.
// fn must not call methods on s.
//
// Use [SyncSet.WithRWLock] when write operations are required.
// Use [SyncSet.Range] when only element iteration is required.
func (s *SyncSet[S, T]) WithRLock(fn func(set.SetLike[T])) {
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
func (s *SyncSet[S, T]) WithLock(fn func(S)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(s.values)
}
