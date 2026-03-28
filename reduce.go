package set

import (
	"iter"
)

// AccumulateTry returns an iterator of intermediate accumulation results.
// Each step applies fn to the current accumulator and the next set.
// If fn returns false, the result of that step is yielded and iteration stops.
// Iteration also stops if the consumer stops early.
// For zero sets, a single zero value is yielded.
// For one set, no values are yielded.
func AccumulateTry[T comparable, S SetLike[T]](
	sets []S,
	fn func(a, b S) (S, bool),
) iter.Seq[S] {
	return func(yield func(S) bool) {
		if len(sets) == 0 {
			var s S
			yield(s)
			return
		}

		r := sets[0]
		for _, s := range sets[1:] {
			var ok bool
			r, ok = fn(r, s)
			if !yield(r) || !ok {
				return
			}
		}
	}
}

// Accumulate returns a iterator of intermediate accumulation results.
// Each step applies fn to the current accumulator and the next set.
// fn must return a new set and must not mutate its arguments.
// Iteration stops if the consumer stops early.
// For zero sets, a single zero value is yielded.
// For one set, no values are yielded.
func Accumulate[T comparable, S SetLike[T]](
	sets []S,
	fn func(a, b S) S,
) iter.Seq[S] {
	return func(yield func(S) bool) {
		if len(sets) == 0 {
			var zero S
			yield(zero)
			return
		}
		r := sets[0]
		for _, s := range sets[1:] {
			r = fn(r, s)
			if !yield(r) {
				return
			}
		}
	}
}

// ReduceTry reduces sets using fn, stopping early if fn returns false.
// The result of the step that signalled stop is returned as the final value.
// If no sets are provided, a zero value is returned.
func ReduceTry[T comparable, S SetLike[T]](
	sets []S,
	fn func(a, b S) (S, bool),
) S {
	if len(sets) == 0 {
		var zero S
		return zero
	}

	r := sets[0]

	for i := 1; i < len(sets); i++ {
		var ok bool
		r, ok = fn(r, sets[i])
		if !ok {
			break
		}
	}
	return r
}

// ReduceWhile reduces sets using fn, stopping when pred returns false for a result.
// The last result that passes pred is returned.
// If no sets are provided, a zero value is returned.
func ReduceWhile[T comparable, S SetLike[T]](
	sets []S,
	fn func(a, b S) S,
	pred func(s S) bool,
) S {
	if len(sets) == 0 {
		var zero S
		return zero
	}

	r := sets[0]

	for i := 1; i < len(sets); i++ {
		r = fn(r, sets[i])
		if !pred(r) {
			return r
		}
	}
	return r
}

// ReduceUntil reduces sets using fn, stopping when pred returns true for a result.
// The first result that passes pred is returned.
// If no sets are provided, a zero value is returned.
func ReduceUntil[T comparable, S SetLike[T]](
	sets []S,
	fn func(a, b S) S,
	pred func(s S) bool,
) S {
	if len(sets) == 0 {
		var zero S
		return zero
	}

	r := sets[0]

	for i := 1; i < len(sets); i++ {
		r = fn(r, sets[i])
		if pred(r) {
			return r
		}
	}
	return r
}

// Reduce reduces sets by applying fn pairwise from left to right.
// The final accumulated value is returned.
// If no sets are provided, a zero value is returned.
func Reduce[T comparable, S SetLike[T]](
	sets []S,
	fn func(a, b S) S,
) S {
	if len(sets) == 0 {
		var zero S
		return zero
	}

	r := sets[0]

	for i := 1; i < len(sets); i++ {
		r = fn(r, sets[i])
	}
	return r
}
