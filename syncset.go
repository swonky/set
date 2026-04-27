package set

import (
	"sync"

	"github.com/swonky/set/types"
)

var (
	_ types.SetLike[int]     = (*SyncSet[int])(nil)
	_ types.MutableSet[int]  = (*SyncSet[int])(nil)
	_ types.LockableSet[int] = (*SyncSet[int])(nil)
)

// SyncSet is a concurrency-safe wrapper around a MutableSet.
//
// All access to the wrapped set is coordinated by an internal RWMutex.
// The wrapped set itself need not provide synchronization.
//
// The zero value is not ready for use. Use [NewSync] or [Wrap].
type SyncSet[T any] struct {
	mu     sync.RWMutex
	values types.MutableSet[T]
}

// Wrap returns a new [SyncSet] wrapping ms.
//
// After wrapping, ms must not be accessed directly unless the caller
// provides external synchronization compatible with the returned SyncSet.
func Wrap[T any](ms types.MutableSet[T]) *SyncSet[T] {
	return &SyncSet[T]{values: ms}
}

// Len returns the number of elements in the set.
func (s *SyncSet[T]) Len() int {
	s.mu.RLock()
	n := s.values.Len()
	s.mu.RUnlock()
	return n
}

// Contains reports whether item is present in the set.
func (s *SyncSet[T]) Contains(item T) bool {
	s.mu.RLock()
	ok := s.values.Contains(item)
	s.mu.RUnlock()
	return ok
}

// Add inserts item into the set.
func (s *SyncSet[T]) Add(item T) {
	s.mu.Lock()
	s.values.Add(item)
	s.mu.Unlock()
}

// Delete removes item from the set.
func (s *SyncSet[T]) Delete(item T) {
	s.mu.Lock()
	s.values.Delete(item)
	s.mu.Unlock()
}

// Range can be used as an [iter.Seq].
//
// It calls yield for each element while holding a read lock.
// Iteration stops when yield returns false.
// Iteration order is implementation-defined.
//
// yield must not call methods on s.
// Use [SyncSet.WithRLock] or [SyncSet.WithLock] when multiple
// operations must be coordinated within one callback.
func (s *SyncSet[T]) Range(yield func(T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.values.Range(func(t T) bool {
		return yield(t)
	})
}

// WithRLock calls fn while holding a read lock and passes a callback-scoped
// read handle to the wrapped set.
//
// The handle must not be retained after fn returns.
// fn must use the provided handle rather than calling methods on s.
func (s *SyncSet[T]) WithRLock(fn func(types.SetLike[T])) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fn(s.values)
}

// WithLock calls fn while holding an exclusive lock and passes a
// callback-scoped writable handle to the wrapped set.
//
// The handle must not be retained after fn returns.
// fn must use the provided handle rather than calling methods on s.
func (s *SyncSet[T]) WithLock(fn func(types.MutableSet[T])) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fn(s.values)
}
