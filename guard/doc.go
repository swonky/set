// Package guard provides coordination helpers for performing multi-step set
// operations against values from package types.
//
// The package wraps one or two sets in callback-based guards that apply the
// requested read or write access mode for the duration of a function call.
// When a set implements types.LockableSet, the corresponding lock is held
// while the callback executes. Sets that do not provide locking are still
// passed through the same callback-scoped API without additional
// synchronization.
//
// For two-set operations, guard acquires locks in a stable order when needed
// and coalesces duplicate references to the same lockable set into a single
// acquisition using the strongest required access mode.
//
// Callback parameters are temporary handles valid only for the duration of
// the callback. Code should use those handles rather than the original set
// values while the callback is executing.
//
// Typical use:
//
//	guard.Write(ms).Do(func(wset types.MutableSet[int]) {
//		wset.Add(1)
//	})
//
//	guard.ReadWrite(lhs, rhs).Do(func(
//		rset types.SetLike[int],
//		wset types.MutableSet[int],
//	) {
//		rset.Range(func(v int) bool {
//			wset.Add(v)
//			return true
//		})
//	})
package guard
