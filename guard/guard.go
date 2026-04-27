package guard

import (
	"reflect"

	"github.com/swonky/set/types"
)

type guard[S, T any] struct {
	s     S
	write bool
}

// Do calls fn with the bound set.
// If the set supports locking, and write is false, a read lock is used.
// If the set supports locking, and write is true, an exclusive lock is used.
func (g *guard[S, T]) Do(fn func(S)) {
	if ls, ok := any(g.s).(types.LockableSet[T]); ok {
		if g.write {
			ls.WithLock(func(ms types.MutableSet[T]) {
				fn(ms.(S))
			})
		} else {
			ls.WithRLock(func(sl types.SetLike[T]) {
				fn(sl.(S))
			})
		}
		return
	}

	fn(g.s)
}

type guard2[S1 types.SetLike[T1], S2 types.SetLike[T2], T1, T2 any] struct {
	lhs *guard[S1, T1]
	rhs *guard[S2, T2]
}

// Do calls fn while holding the appropriate lock and passes a temporary view of each set for use only during the callback.
//
// If both sets are lockable, locks are acquired in a stable order to avoid deadlocks.
// If both guards reference the same underlying set, a single lock is acquired using the strongest required access mode.
//
// lhs and rhs must not be retained after fn returns.
// fn must not call methods on s.
func (g *guard2[S1, S2, T1, T2]) Do(fn func(lhs S1, rhs S2)) {
	id1 := identity(g.lhs.s)
	id2 := identity(g.rhs.s)

	switch {
	case id1 == id2 && id1 != 0:
		copylhs := g.lhs
		copylhs.write = copylhs.write || g.rhs.write

		copylhs.Do(func(s1 S1) {
			fn(s1, any(s1).(S2))
		})

	case id1 < id2:
		g.lhs.Do(func(s1 S1) {
			g.rhs.Do(func(s2 S2) {
				fn(s1, s2)
			})
		})

	default:
		g.rhs.Do(func(s2 S2) {
			g.lhs.Do(func(s1 S1) {
				fn(s1, s2)
			})
		})
	}
}

// newGuard2 returns a guard2.
// This exists to ease the generic-typed pain in my head.
func newGuard2[
	S1 types.SetLike[T1],
	S2 types.SetLike[T2],
	T1, T2 any,
](
	lhs *guard[S1, T1],
	rhs *guard[S2, T2],
) *guard2[S1, S2, T1, T2] {
	return &guard2[S1, S2, T1, T2]{
		lhs: lhs,
		rhs: rhs,
	}
}

// read returns a guard that provides read access to s.
func read[T any](s types.SetLike[T]) *guard[types.SetLike[T], T] {
	return &guard[types.SetLike[T], T]{
		s:     s,
		write: false,
	}
}

// write returns a guard that provides write access to ms.
func write[T any](ms types.MutableSet[T]) *guard[types.MutableSet[T], T] {
	return &guard[types.MutableSet[T], T]{
		s:     ms,
		write: true,
	}
}

// identity returns the underlying pointer identity of lockable pointer-backed sets.
// Non-lockable or non-pointer-backed sets return zero.
func identity[S types.SetLike[T], T any](s S) uintptr {
	if ls, ok := any(s).(types.LockableSet[T]); ok {
		rv := reflect.ValueOf(ls)
		if rv.Kind() == reflect.Pointer {
			return rv.Pointer()
		}
	}
	return 0
}

// # Public functions

// SyncR returns a guard that provides read access to s.
func Read[T any](s types.SetLike[T]) SyncReader[T] {
	return read(s)
}

// SyncR returns a guard that provides read access to s.
func Write[T any](ms types.MutableSet[T]) SyncWriter[T] {
	return write(ms)
}

// ReadRead returns a two-set guard with read access to both sets.
func ReadRead[T1, T2 any](lhs types.SetLike[T1], rhs types.SetLike[T2]) SyncReadRead[T1, T2] {
	return newGuard2(read(lhs), read(rhs))
}

// ReadWrite returns a two-set guard with read access to the first set and
// write access to the second set.
func ReadWrite[T1, T2 any](lhs types.SetLike[T1], rhs types.MutableSet[T2]) SyncReadWrite[T1, T2] {
	return newGuard2(read(lhs), write(rhs))
}

// WriteRead returns a two-set guard with read access to the first set and
// write access to the second set.
func WriteRead[T1, T2 any](lhs types.MutableSet[T1], rhs types.SetLike[T2]) SyncWriteRead[T1, T2] {
	return newGuard2(write(lhs), read(rhs))
}

// WriteWrite returns a two-set guard with write access to both sets.
func WriteWrite[T1, T2 any](lhs types.MutableSet[T1], rhs types.MutableSet[T2]) SyncWriteWrite[T1, T2] {
	return newGuard2(write(lhs), write(rhs))
}

// Same reports whether lhs and rhs refer to the same underlying lockable set.
// It returns true only when both values expose the same stable lock identity.
// Non-lockable values, nil values, or values without a comparable lock
// identity report false.
func Same[T1, T2 any](lhs types.SetLike[T1], rhs types.SetLike[T2]) bool {
	return identity(lhs) != 0 && identity(lhs) == identity(rhs)
}
