package set

import (
	"cmp"
	"iter"
	"slices"
)

// AddCheck inserts an item into the set. Returns true if item was already present.
func AddCheck[T any](ms MutableSet[T], elem T) bool {
	if ms.Contains(elem) {
		return true
	}
	ms.Add(elem)
	return false
}

func Add[T any](ms MutableSet[T], elems ...T) {
	for _, v := range elems {
		ms.Add(v)
	}
}

func AddDedup[S MutableSet[T], T cmp.Ordered](s S, elems ...T) {
	if len(elems) == 0 {
		return
	}

	slices.Sort(elems)
	elems = slices.Compact(elems)

	for _, v := range elems {
		s.Add(v)
	}
}

func AsYielder[T any](fn func(T)) func(T) bool {
	return func(t T) bool {
		fn(t)
		return true
	}
}

func Append[T any](ms MutableSet[T], elems ...T) {
	slices.Values(elems)(AsYielder(ms.Add))
}

func Copy[S1 MutableSet[T], S2 SetLike[T], T comparable](dst S1, src S2) {
	src.Range(func(elem T) bool {
		dst.Add(elem)
		return true
	})
}

func Consume[MS MutableSet[T], T any](ms MS) iter.Seq[T] {
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

func Pop[MS MutableSet[T], T any](ms MS) (T, bool) {
	type result[T any] struct {
		v  T
		ok bool
	}
	out := WithWrite(ms, func(sl MS) result[T] {
		var r result[T]
		for t := range sl.Range {
			r.v = t
			r.ok = true
			sl.Delete(t)
			return r
		}
		return r
	})
	return out.v, out.ok
}

func Clear[MS MutableSet[T], T any](ms MS) {
	if ms.Len() == 0 {
		return
	}
	WithWrite(ms, func(ms2 MS) struct{} {
		ms2.Range(func(t T) bool {
			ms2.Delete(t)
			return true
		})
		return struct{}{}
	})
}

func UnionInto[T any](dst MutableSet[T], sets ...SetLike[T]) {
	for _, src := range sets {
		op2(src, dst, func(a, b SetLike[T]) struct{} {
			ms := b.(MutableSet[T])
			a.Range(func(t T) bool {
				ms.Add(t)
				return true
			})
			return struct{}{}
		})
	}
}

func IntersectWith[T any](dst MutableSet[T], src SetLike[T]) {
	op2(dst, src, func(a, b SetLike[T]) struct{} {
		ms := a.(MutableSet[T])

		a.Range(func(t T) bool {
			if !b.Contains(t) {
				ms.Delete(t)
			}
			return true
		})
		return struct{}{}
	})
}

func Diff[T any](dst MutableSet[T], src SetLike[T]) {
	op2(dst, src, func(a, b SetLike[T]) struct{} {
		md := a.(MutableSet[T])
		for t := range b.Range {
			if md.Contains(t) {
				dst.Delete(t)
			}
		}
		return struct{}{}
	})
}

func SymDiff[T any](dst MutableSet[T], src SetLike[T]) {
	op2(dst, src, func(a, b SetLike[T]) struct{} {
		md := a.(MutableSet[T])
		for t := range b.Range {
			if md.Contains(t) {
				dst.Delete(t)
			} else {
				dst.Add(t)
			}
		}
		return struct{}{}
	})
}
