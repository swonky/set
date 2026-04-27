package types

// SetLike is the common read-only interface implemented by set types.
type SetLike[T any] interface {
	// Contains reports whether item is present.
	Contains(item T) bool

	// Range calls yield for each element until all elements have been visited or yield returns false.
	// Iteration order is implementation-defined.
	//
	// The method value Range has type [iter.Seq].
	Range(yield func(T) bool)

	// Len returns the number of elements currently in the set.
	Len() int
}

// MutableSet is a set that permits mutation.
type MutableSet[T any] interface {
	SetLike[T]

	// Range calls yield for each element until all elements have been visited or yield returns false.
	// Iteration order is implementation-defined.
	//
	// The method value Range has type [iter.Seq].
	//
	// Implementations of [MutableSet] must permit [MutableSet.Add] and [MutableSet.Delete] on the receiver during iteration.
	// Elements added during iteration may be visited or skipped.
	// Elements deleted before being visited must not be visited afterwards.
	Range(yield func(T) bool)

	// Add inserts elem into the set.
	Add(elem T)

	// Delete removes elem from the set if present.
	Delete(elem T)
}

type Clearable[T any] interface {
	MutableSet[T]

	// Clear removes all elements from the set
	Clear()
}

// LockableSet is a set that exposes coordinated read and write access.
type LockableSet[T any] interface {
	SetLike[T]

	// WithRLock calls fn while a shared read lock is held.
	// The callback must not retain rset after fn returns.
	WithRLock(func(rset SetLike[T]))

	// WithLock calls fn while an exclusive write lock is held.
	// The callback must not retain wset after fn returns.
	WithLock(func(wset MutableSet[T]))
}
