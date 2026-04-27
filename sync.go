package set

// type syncR[T any] struct {
// 	s types.SetLike[T]
// }

// type syncW[MS types.MutableSet[T], T any] struct {
// 	s MS
// }

// func (x syncR[T]) R(fn func(types.SetLike[T])) {
// 	if l, ok := x.s.(LockableSet[types.MutableSet[T], T]); ok {
// 		l.WithRLock(fn)
// 		return
// 	}
// 	fn(x.s)
// }

// func (x syncW[MS, T]) W(fn func(MS)) {
// 	if l, ok := x.s.(LockableSet[MS, T]); ok {
// 		l.WithLock(fn)
// 		return
// 	}

// 	fn(m)
// }

// func SyncReader[T any](s types.SetLike[T]) func(fn func(types.SetLike[T])) {
// 	return func(fn func(types.SetLike[T])) {
// 		if l, ok := s.(types.LockableSet[T]); ok {
// 			l.WithRLock(fn)
// 			return
// 		}
// 		fn(s)
// 	}
// }

// func SyncWriter[MS types.MutableSet[T], T any](s types.MutableSet[T]) func(fn func(types.MutableSet[T])) {
// 	return func(fn func(types.MutableSet[T])) {
// 		if l, ok := s.(types.LockableSet[T]); ok {
// 			l.WithLock(fn)
// 			return
// 		}
// 		fn(s)
// 	}
// }

// func (sp *setPair[S1, S2, T]) Run(fn func(S1, S2)) {
// 	fn(sp.lhs, sp.rhs)
// }

// func SyncRR[T any](s1 types.SetLike[T], s2 types.SetLike[T]) setPair[types.SetLike[T], types.SetLike[T], T] {
// 	return setPair[types.SetLike[T], types.SetLike[T], T]{lhs: s1, rhs: s2}
// }

// func SyncRW[T any](s types.SetLike[T], ms types.MutableSet[T]) setPair[types.SetLike[T], types.MutableSet[T], T] {
// 	return setPair[types.SetLike[T], types.MutableSet[T], T]{lhs: s, rhs: ms}
// }

// func SyncWW[T any](ms1 types.MutableSet[T], ms2 types.MutableSet[T]) setPair[types.MutableSet[T], types.MutableSet[T], T] {
// 	return setPair[types.MutableSet[T], types.MutableSet[T], T]{lhs: ms1, rhs: ms2}
// }

// // type SyncRR[T any] setPair[types.SetLike[T], types.SetLike[T], T]
// // type SyncRW[T any] setPair[types.SetLike[T], types.MutableSet[T], T]
// // type SyncWR[T any] setPair[types.MutableSet[T], types.SetLike[T], T]
// // type SyncWW[T any] setPair[types.MutableSet[T], types.MutableSet[T], T]

// // func (sp *setPair[S1,S2,T]) R

// // func RR(s types.SetLike[T]) setPair[types.SetLike[T], types.SetLike[T]] {

// // }

// func (p setPair[T]) add(s types.SetLike[T]) setPair[T] {
// 	if p.lhs == nil {
// 		p.lhs = s
// 		return p
// 	}
// 	p.rhs = s
// 	return p
// }

// type PairBuilderA[T any] interface {
// 	WithRead1(types.SetLike[T]) PairBuilderB[T]
// 	WithWrite1(types.MutableSet[T]) PairBuilderB[T]
// }

// type PairBuilderB[T any] interface {
// 	WithRead2(types.SetLike[T])
// 	WithWrite2(types.MutableSet[T])
// }

// func (p setPair[T]) WithRead1(s types.SetLike[T]) PairBuilderB[T] { return p.add(s) }

// func (p setPair[T]) WithWrite1(s types.MutableSet[T]) PairBuilderB[T] { return p.add(s) }

// func (p setPair[T]) WithRead2(s types.SetLike[T]) { p.add(s) }

// func (p setPair[T]) WithWrite2(s types.MutableSet[T]) { p.add(s) }

// // sync1 coordinates synchronized access to a single set.
// type sync1[T any] struct {
// 	s types.SetLike[T]
// }

// // sync2 coordinates synchronized access to two sets.
// type sync2[A, B any] struct {
// 	a SetLike[A]
// 	b SetLike[B]
// }

// // Sync binds one set and returns a synchronization coordinator.
// func Sync[T any](s types.SetLike[T]) sync1[T] {
// 	return sync1[T]{s: s}
// }

// // Sync2 binds two sets and returns a synchronization coordinator.
// func Sync2[A, B any](a SetLike[A], b SetLike[B]) sync2[A, B] {
// 	return sync2[A, B]{a: a, b: b}
// }

// // R calls fn with read access to the bound set.
// // If the set supports locking, a read lock is held for the duration
// // of the callback.
// func (x sync1[T]) R(fn func(types.SetLike[T])) {
// 	if l, ok := x.s.(LockableSet[types.MutableSet[T], T]); ok {
// 		l.WithRLock(fn)
// 		return
// 	}
// 	fn(x.s)
// }

// // W calls fn with write access to the bound set.
// // If the set supports locking, an exclusive lock is held for the
// // duration of the callback.
// func (x sync1[T]) W(fn func(types.MutableSet[T])) {
// 	m := x.s.(types.MutableSet[T])

// 	if l, ok := x.s.(LockableSet[types.MutableSet[T], T]); ok {
// 		l.WithLock(fn)
// 		return
// 	}

// 	fn(m)
// }

// type SetWrapper[S types.SetLike[T], T any] struct {
// 	s S
// }

// func (s SetWrapper[S, T]) Contains(item T) bool     { return s.s.Contains(item) }
// func (s SetWrapper[S, T]) Range(yield func(T) bool) { s.s.Range(yield) }
// func (s SetWrapper[S, T]) Len() int                 { return s.s.Len() }

// func (s SetWrapper[S, T]) WithRLock(fn func(types.SetLike[T])) {
// 	if ls, ok := types.SetLike[T](s).(types.LockableSet[T]); ok {
// 		ls.WithRLock(fn)
// 	}
// 	fn(s)
// }

// func (s SetWrapper[S, T]) WithLock(fn func(types.MutableSet[T])) int {
// 	if ls, ok := types.SetLike[T](s).(types.LockableSet[T]); ok {
// 		ls.WithRLock(fn)
// 	}
// 	fn(s)
// }

// // RR calls fn with read access to both bound sets.
// // When both sets support locking, locks are acquired in a stable order
// // to avoid deadlocks.
// func RR[S1, S2 types.SetLike[T], T any](s1 S1, s2 S2, fn func(sl1 types.SetLike[T], sl2 types.SetLike[T])) {
// 	la, oka := types.SetLike[T](s1).(types.LockableSet[T])
// 	lb, okb := types.SetLike[T](s2).(types.LockableSet[T])

// 	run := func() { fn(s1, s2) }

// 	switch {
// 	case !oka && !okb:
// 		run()

// 	case oka && okb && samePtr(la, lb):
// 		la.WithRLock(func(types.SetLike[T]) {
// 			run()
// 		})

// 	case oka && okb:
// 		if ptrOf(la) < ptrOf(lb) {
// 			la.WithRLock(func(SetLike[A]) {
// 				lb.WithRLock(func(SetLike[B]) {
// 					run()
// 				})
// 			})
// 		} else {
// 			lb.WithRLock(func(SetLike[B]) {
// 				la.WithRLock(func(SetLike[A]) {
// 					run()
// 				})
// 			})
// 		}

// 	case oka:
// 		la.WithRLock(func(SetLike[A]) {
// 			run()
// 		})

// 	default:
// 		lb.WithRLock(func(SetLike[B]) {
// 			run()
// 		})
// 	}
// }

// func orderPtr[T any](a T, b T) (T, T) {

// }

// func RR2[T any](s1, s2 types.SetLike[T], fn func(sl1 types.SetLike[T], sl2 types.SetLike[T])) {

// 	if la, ok := s1.(types.LockableSet[T]); ok {
// 		if lb, ok := s2.(types.LockableSet[T]); ok {
// 			la.WithRLock(func(a2 types.SetLike[T]) {
// 				lb.WithRLock(func(b2 types.SetLike[T]) {
// 					fn(a2, b2)
// 				})
// 			})
// 		}
// 	}
// 	fn(p.lhs, p.rhs)
// }

// // RR calls fn with read access to both bound sets.
// // When both sets support locking, locks are acquired in a stable order
// // to avoid deadlocks.
// func RR[S1, S2 types.SetLike[T], T any](s1 S1, s2 S2, fn func(sl1 types.SetLike[T], sl2 types.SetLike[T])) {
// 	ls1, oka := types.SetLike[T](s1).(types.LockableSet[T])
// 	ls2, okb := types.SetLike[T](s2).(types.LockableSet[T])

// 	switch {
// 	case !oka && !okb:
// 		fn(s1, s2)

// 	case oka && okb && samePtr(ls1, ls2):
// 		ls1.WithRLock(func(lockedS1 types.SetLike[T]) { fn(lockedS1, lockedS1) })

// 	case oka && okb:
// 		if ptrOf(ls1) < ptrOf(ls2) {
// 			ls1.WithRLock(func(lockedS1 types.SetLike[T]) {
// 				ls2.WithRLock(func(lockedS2 types.SetLike[T]) { fn(lockedS1, lockedS2) })
// 			})
// 		} else {
// 			ls2.WithRLock(func(lockedS2 types.SetLike[T]) {
// 				ls1.WithRLock(func(lockedS1 types.SetLike[T]) { fn(lockedS1, lockedS2) })
// 			})
// 		}

// 	case oka:
// 		ls1.WithRLock(func(lockedS1 types.SetLike[T]) { fn(lockedS1, s2) })

// 	default:
// 		ls2.WithRLock(func(lockedS2 types.SetLike[T]) { fn(s1, lockedS2) })
// 	}
// }

// // RW calls fn with read access to the first bound set and write access
// // to the second bound set. When both sets support locking, locks are
// // acquired in a stable order to avoid deadlocks.
// func (x sync2[A, B]) RW(fn func(s SetLike[A], ms MutableSet[B])) {
// 	mb := x.b.(MutableSet[B])

// 	la, oka := x.a.(LockableSet[MutableSet[A], A])
// 	lb, okb := x.b.(LockableSet[MutableSet[B], B])

// 	run := func() {
// 		fn(x.a, mb)
// 	}

// 	switch {
// 	case !oka && !okb:
// 		run()

// 	case oka && okb && samePtr(la, lb):
// 		lb.WithLock(func(MutableSet[B]) {
// 			run()
// 		})

// 	case oka && okb:
// 		if ptrOf(la) < ptrOf(lb) {
// 			la.WithRLock(func(SetLike[A]) {
// 				lb.WithLock(func(MutableSet[B]) {
// 					run()
// 				})
// 			})
// 		} else {
// 			lb.WithLock(func(MutableSet[B]) {
// 				la.WithRLock(func(SetLike[A]) {
// 					run()
// 				})
// 			})
// 		}

// 	case oka:
// 		la.WithRLock(func(SetLike[A]) {
// 			run()
// 		})

// 	default:
// 		lb.WithLock(func(MutableSet[B]) {
// 			run()
// 		})
// 	}
// }

// // WR calls fn with write access to the first bound set and read access
// // to the second bound set. When both sets support locking, locks are
// // acquired in a stable order to avoid deadlocks.
// func (x sync2[A, B]) WR(fn func(ms MutableSet[A], s SetLike[B])) {
// 	ma := x.a.(MutableSet[A])

// 	la, oka := x.a.(LockableSet[MutableSet[A], A])
// 	lb, okb := x.b.(LockableSet[MutableSet[B], B])

// 	run := func() {
// 		fn(ma, x.b)
// 	}

// 	switch {
// 	case !oka && !okb:
// 		run()

// 	case oka && okb && samePtr(la, lb):
// 		la.WithLock(func(MutableSet[A]) {
// 			run()
// 		})

// 	case oka && okb:
// 		if ptrOf(la) < ptrOf(lb) {
// 			la.WithLock(func(MutableSet[A]) {
// 				lb.WithRLock(func(SetLike[B]) {
// 					run()
// 				})
// 			})
// 		} else {
// 			lb.WithRLock(func(SetLike[B]) {
// 				la.WithLock(func(MutableSet[A]) {
// 					run()
// 				})
// 			})
// 		}

// 	case oka:
// 		la.WithLock(func(MutableSet[A]) {
// 			run()
// 		})

// 	default:
// 		lb.WithRLock(func(SetLike[B]) {
// 			run()
// 		})
// 	}
// }

// // WW calls fn with write access to both bound sets.
// // When both sets support locking, locks are acquired in a stable order
// // to avoid deadlocks.
// func (x sync2[A, B]) WW(fn func(ms1 MutableSet[A], ms2 MutableSet[B])) {
// 	ma := x.a.(MutableSet[A])
// 	mb := x.b.(MutableSet[B])

// 	la, oka := x.a.(LockableSet[MutableSet[A], A])
// 	lb, okb := x.b.(LockableSet[MutableSet[B], B])

// 	run := func() {
// 		fn(ma, mb)
// 	}

// 	switch {
// 	case !oka && !okb:
// 		run()

// 	case oka && okb && samePtr(la, lb):
// 		la.WithLock(func(MutableSet[A]) {
// 			run()
// 		})

// 	case oka && okb:
// 		if ptrOf(la) < ptrOf(lb) {
// 			la.WithLock(func(MutableSet[A]) {
// 				lb.WithLock(func(MutableSet[B]) {
// 					run()
// 				})
// 			})
// 		} else {
// 			lb.WithLock(func(MutableSet[B]) {
// 				la.WithLock(func(MutableSet[A]) {
// 					run()
// 				})
// 			})
// 		}

// 	case oka:
// 		la.WithLock(func(MutableSet[A]) {
// 			run()
// 		})

// 	default:
// 		lb.WithLock(func(MutableSet[B]) {
// 			run()
// 		})
// 	}
// }

// // ptrOf returns an identity value used for lock ordering.
// func ptrOf[T any](v types.LockableSet[T]) uintptr {
// 	return uintptr(unsafe.Pointer(&v))
// }

// // samePtr reports whether a and b share the same identity.
// func samePtr[A, B any](a LockableSet[A], b LockableSet[B]) bool {
// 	return ptrOf(a) == ptrOf(b)
// }
