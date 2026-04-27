package guard

import "github.com/swonky/set/types"

var (
	_ syncGuard[types.SetLike[int], int]                                       = SyncReader[int](nil)
	_ syncGuard[types.MutableSet[int], int]                                    = SyncWriter[int](nil)
	_ syncGuard2[types.SetLike[int8], types.SetLike[int16], int8, int16]       = SyncReadRead[int8, int16](nil)
	_ syncGuard2[types.MutableSet[int8], types.MutableSet[int16], int8, int16] = SyncWriteWrite[int8, int16](nil)
	_ syncGuard2[types.SetLike[int8], types.MutableSet[int16], int8, int16]    = SyncReadWrite[int8, int16](nil)
	_ syncGuard2[types.MutableSet[int8], types.SetLike[int16], int8, int16]    = SyncWriteRead[int8, int16](nil)
)

// SyncReader coordinates synchronized read access to one set.
type SyncReader[T any] interface {
	// Do calls fn with a callback-scoped read handle to the bound set.
	//
	// If supported by the set, the a read sync lock is held for
	// the duration of fn.
	//
	// rset must not be retained after fn returns.
	// fn must use rset rather than the original set value.
	Do(fn func(rset types.SetLike[T]))
}

// SyncWriter coordinates synchronized write access to one set.
type SyncWriter[T any] interface {
	// Do calls fn with a callback-scoped writable handle to the bound set.
	//
	// If supported by the set, an exclusive sync lock is held for
	// the duration of fn.
	//
	// wset must not be retained after fn returns.
	// fn must use wset rather than the original set value.
	Do(fn func(wset types.MutableSet[T]))
}

// SyncReadRead coordinates synchronized read access to two sets.
type SyncReadRead[T1, T2 any] interface {
	// Do calls fn with callback-scoped read handles to the bound sets.
	//
	// If supported by the sets, shared read access is held for the
	// duration of fn.
	//
	// lhs and rhs must not be retained after fn returns.
	// fn must use lhs and rhs rather than the original set values.
	Do(fn func(lhs types.SetLike[T1], rhs types.SetLike[T2]))
}

// SyncWriteWrite coordinates synchronized write access to two sets.
type SyncWriteWrite[T1, T2 any] interface {
	// Do calls fn with callback-scoped writable handles to the bound sets.
	//
	// If supported by the sets, exclusive write access is held for the
	// duration of fn.
	//
	// lhs and rhs must not be retained after fn returns.
	// fn must use lhs and rhs rather than the original set values.
	Do(fn func(lhs types.MutableSet[T1], rhs types.MutableSet[T2]))
}

// SyncReadWrite coordinates synchronized read access to the lhs set, and
// write access to the rhs set.
type SyncReadWrite[T1, T2 any] interface {
	// Do calls fn with a callback-scoped read handle to the lhs set and a
	// callback-scoped writable handle to the rhs set.
	//
	// If supported by the sets, shared read access is held for lhs and
	// exclusive write access is held for rhs for the duration of fn.
	//
	// lhs and rhs must not be retained after fn returns.
	// fn must use lhs and rhs rather than the original set values.
	Do(fn func(lhs types.SetLike[T1], rhs types.MutableSet[T2]))
}

// SyncWriteRead coordinates synchronized write access to the lhs set, and
// read access to the rhs set.
type SyncWriteRead[T1, T2 any] interface {
	// Do calls fn with a callback-scoped writable handle to the lhs set and
	// a callback-scoped read handle to the rhs set.
	//
	// If supported by the sets, exclusive write access is held for lhs and
	// shared read access is held for rhs for the duration of fn.
	//
	// lhs and rhs must not be retained after fn returns.
	// fn must use lhs and rhs rather than the original set values.
	Do(fn func(lhs types.MutableSet[T1], rhs types.SetLike[T2]))
}
