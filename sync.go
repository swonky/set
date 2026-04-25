package set

import "unsafe"

type syncR[T any] struct {
	s SetLike[T]
}

type syncW[MS MutableSet[T], T any] struct {
	s MS
}

func (x syncR[T]) R(fn func(SetLike[T])) {
	if l, ok := x.s.(LockableSet[MutableSet[T], T]); ok {
		l.WithRLock(fn)
		return
	}
	fn(x.s)
}

func (x syncW[MS, T]) W(fn func(MS)) {
}

func SyncReader[T any](s SetLike[T]) syncR[T] {
	return syncR[T]{s}
}

// sync1 coordinates synchronized access to a single set.
type sync1[T any] struct {
	s SetLike[T]
}

// sync2 coordinates synchronized access to two sets.
type sync2[A, B any] struct {
	a SetLike[A]
	b SetLike[B]
}

// Sync binds one set and returns a synchronization coordinator.
func Sync[T any](s SetLike[T]) sync1[T] {
	return sync1[T]{s: s}
}

// Sync2 binds two sets and returns a synchronization coordinator.
func Sync2[A, B any](a SetLike[A], b SetLike[B]) sync2[A, B] {
	return sync2[A, B]{a: a, b: b}
}

// R calls fn with read access to the bound set.
// If the set supports locking, a read lock is held for the duration
// of the callback.
func (x sync1[T]) R(fn func(SetLike[T])) {
	if l, ok := x.s.(LockableSet[MutableSet[T], T]); ok {
		l.WithRLock(fn)
		return
	}
	fn(x.s)
}

// W calls fn with write access to the bound set.
// If the set supports locking, an exclusive lock is held for the
// duration of the callback.
func (x sync1[T]) W(fn func(MutableSet[T])) {
	m := x.s.(MutableSet[T])

	if l, ok := x.s.(LockableSet[MutableSet[T], T]); ok {
		l.WithLock(fn)
		return
	}

	fn(m)
}

// RR calls fn with read access to both bound sets.
// When both sets support locking, locks are acquired in a stable order
// to avoid deadlocks.
func (x sync2[A, B]) RR(fn func(s1 SetLike[A], s2 SetLike[B])) {
	la, oka := x.a.(LockableSet[MutableSet[A], A])
	lb, okb := x.b.(LockableSet[MutableSet[B], B])

	run := func() {
		fn(x.a, x.b)
	}

	switch {
	case !oka && !okb:
		run()

	case oka && okb && samePtr(la, lb):
		la.WithRLock(func(SetLike[A]) {
			run()
		})

	case oka && okb:
		if ptrOf(la) < ptrOf(lb) {
			la.WithRLock(func(SetLike[A]) {
				lb.WithRLock(func(SetLike[B]) {
					run()
				})
			})
		} else {
			lb.WithRLock(func(SetLike[B]) {
				la.WithRLock(func(SetLike[A]) {
					run()
				})
			})
		}

	case oka:
		la.WithRLock(func(SetLike[A]) {
			run()
		})

	default:
		lb.WithRLock(func(SetLike[B]) {
			run()
		})
	}
}

// RW calls fn with read access to the first bound set and write access
// to the second bound set. When both sets support locking, locks are
// acquired in a stable order to avoid deadlocks.
func (x sync2[A, B]) RW(fn func(s SetLike[A], ms MutableSet[B])) {
	mb := x.b.(MutableSet[B])

	la, oka := x.a.(LockableSet[MutableSet[A], A])
	lb, okb := x.b.(LockableSet[MutableSet[B], B])

	run := func() {
		fn(x.a, mb)
	}

	switch {
	case !oka && !okb:
		run()

	case oka && okb && samePtr(la, lb):
		lb.WithLock(func(MutableSet[B]) {
			run()
		})

	case oka && okb:
		if ptrOf(la) < ptrOf(lb) {
			la.WithRLock(func(SetLike[A]) {
				lb.WithLock(func(MutableSet[B]) {
					run()
				})
			})
		} else {
			lb.WithLock(func(MutableSet[B]) {
				la.WithRLock(func(SetLike[A]) {
					run()
				})
			})
		}

	case oka:
		la.WithRLock(func(SetLike[A]) {
			run()
		})

	default:
		lb.WithLock(func(MutableSet[B]) {
			run()
		})
	}
}

// WR calls fn with write access to the first bound set and read access
// to the second bound set. When both sets support locking, locks are
// acquired in a stable order to avoid deadlocks.
func (x sync2[A, B]) WR(fn func(ms MutableSet[A], s SetLike[B])) {
	ma := x.a.(MutableSet[A])

	la, oka := x.a.(LockableSet[MutableSet[A], A])
	lb, okb := x.b.(LockableSet[MutableSet[B], B])

	run := func() {
		fn(ma, x.b)
	}

	switch {
	case !oka && !okb:
		run()

	case oka && okb && samePtr(la, lb):
		la.WithLock(func(MutableSet[A]) {
			run()
		})

	case oka && okb:
		if ptrOf(la) < ptrOf(lb) {
			la.WithLock(func(MutableSet[A]) {
				lb.WithRLock(func(SetLike[B]) {
					run()
				})
			})
		} else {
			lb.WithRLock(func(SetLike[B]) {
				la.WithLock(func(MutableSet[A]) {
					run()
				})
			})
		}

	case oka:
		la.WithLock(func(MutableSet[A]) {
			run()
		})

	default:
		lb.WithRLock(func(SetLike[B]) {
			run()
		})
	}
}

// WW calls fn with write access to both bound sets.
// When both sets support locking, locks are acquired in a stable order
// to avoid deadlocks.
func (x sync2[A, B]) WW(fn func(ms1 MutableSet[A], ms2 MutableSet[B])) {
	ma := x.a.(MutableSet[A])
	mb := x.b.(MutableSet[B])

	la, oka := x.a.(LockableSet[MutableSet[A], A])
	lb, okb := x.b.(LockableSet[MutableSet[B], B])

	run := func() {
		fn(ma, mb)
	}

	switch {
	case !oka && !okb:
		run()

	case oka && okb && samePtr(la, lb):
		la.WithLock(func(MutableSet[A]) {
			run()
		})

	case oka && okb:
		if ptrOf(la) < ptrOf(lb) {
			la.WithLock(func(MutableSet[A]) {
				lb.WithLock(func(MutableSet[B]) {
					run()
				})
			})
		} else {
			lb.WithLock(func(MutableSet[B]) {
				la.WithLock(func(MutableSet[A]) {
					run()
				})
			})
		}

	case oka:
		la.WithLock(func(MutableSet[A]) {
			run()
		})

	default:
		lb.WithLock(func(MutableSet[B]) {
			run()
		})
	}
}

// ptrOf returns an identity value used for lock ordering.
func ptrOf[T any](v LockableSet[MutableSet[T], T]) uintptr {
	return uintptr(unsafe.Pointer(&v))
}

// samePtr reports whether a and b share the same identity.
func samePtr[A, B any](a LockableSet[MutableSet[A], A], b LockableSet[MutableSet[B], B]) bool {
	return ptrOf(a) == ptrOf(b)
}
