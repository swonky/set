package set

import (
	"iter"

	"github.com/swonky/set/guard"
	"github.com/swonky/set/types"
)

// Add inserts elems into ms.
func Add[MS types.MutableSet[T], T any](ms MS, elems ...T) {
	switch len(elems) {
	case 0:
	case 1:
		ms.Add(elems[0])
	default:
		guard.Write(ms).Do(
			func(wset types.MutableSet[T]) {
				for _, v := range elems {
					wset.Add(v)
				}
			})
	}
}

// Consume removes and yields elements from ms until empty.
// Elements are yielded in implementation-defined order.
func Consume[MS types.MutableSet[T], T any](ms MS) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			v, ok := Pop(ms)
			if !ok {
				return
			}

			if !yield(v) {
				return
			}
		}
	}
}

// Copy inserts all elements of src into dst.
// To unite multiple source sets, use [UnionInto].
func Copy[S1 types.MutableSet[T], S2 types.SetLike[T], T comparable](dst S1, src S2) {
	guard.WriteRead(dst, src).Do(
		func(lhs types.MutableSet[T], rhs types.SetLike[T]) {
			rhs.Range(func(elem T) bool {
				lhs.Add(elem)
				return true
			})
		})
}

// Pop removes and returns an arbitrary element from ms.
// If ms is empty, Pop returns the zero value of T and false.
func Pop[MS types.MutableSet[T], T any](ms MS) (T, bool) {
	var (
		ok     bool
		result T
	)

	guard.Write(ms).Do(
		func(wset types.MutableSet[T]) {
			wset.Range(func(t T) bool {
				result = t
				ok = true
				wset.Delete(t)
				return false
			})
		})

	return result, ok
}

// Clear removes all elements from ms.
func Clear[MS types.MutableSet[T], T any](ms MS) {
	if e, ok := any(ms).(types.Clearable[T]); ok {
		e.Clear()
		return
	}

	guard.Write(ms).Do(
		func(wset types.MutableSet[T]) {
			if wset.Len() == 0 {
				return
			}

			wset.Range(func(t T) bool {
				wset.Delete(t)
				return true
			})
		})
}

// UnionInto inserts into dst every element present in srcs.
func UnionInto[T any](dst types.MutableSet[T], srcs ...types.SetLike[T]) {
	for _, src := range srcs {
		guard.WriteRead(dst, src).Do(
			func(
				lhs types.MutableSet[T],
				rhs types.SetLike[T],
			) {
				rhs.Range(func(t T) bool {
					lhs.Add(t)
					return true
				})
			})
	}
}

// IntersectWith removes from dst any element not present in src.
func IntersectWith[T any](dst types.MutableSet[T], src types.SetLike[T]) {
	guard.WriteRead(dst, src).Do(
		func(
			lhs types.MutableSet[T],
			rhs types.SetLike[T],
		) {
			lhs.Range(func(t T) bool {
				if !rhs.Contains(t) {
					lhs.Delete(t)
				}
				return true
			})
		})
}

// Diff removes from dst any element present in src.
func Diff[T any](dst types.MutableSet[T], src types.SetLike[T]) {
	guard.WriteRead(dst, src).Do(
		func(
			lhs types.MutableSet[T],
			rhs types.SetLike[T],
		) {
			lhs.Range(func(t T) bool {
				if rhs.Contains(t) {
					lhs.Delete(t)
				}
				return true
			})
		})
}

// SymDiff updates dst to contain elements present in exactly one of dst or src.
func SymDiff[T any](dst types.MutableSet[T], src types.SetLike[T]) {
	guard.WriteRead(dst, src).Do(
		func(
			lhs types.MutableSet[T],
			rhs types.SetLike[T],
		) {
			rhs.Range(func(t T) bool {
				if lhs.Contains(t) {
					lhs.Delete(t)
				} else {
					lhs.Add(t)
				}
				return true
			})
		})
}
