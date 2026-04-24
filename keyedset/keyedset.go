package keyedset

import (
	"maps"

	"github.com/swonky/set"
	"github.com/swonky/set/internal/base"
)

var _ set.SetLike[Keyed] = KeyedSet[Keyed]{}

type Keyed interface{ Key() uint64 }

type KeyedSet[T any] struct {
	smap map[uint64]T
	fn   func(T) uint64
}

func defaultKeyFunc[T Keyed](t T) uint64 {
	return t.Key()
}

// -- Constructors --

func WithCustom[T any](fn func(T) uint64, cap ...int) KeyedSet[T] {
	return KeyedSet[T]{
		smap: make(map[uint64]T, base.GetCap(cap...)),
		fn:   fn,
	}
}

func FromWithCustom[T any](fn func(T) uint64, elems []T) KeyedSet[T] {
	m := WithCustom(fn, len(elems))
	for _, v := range elems {
		m.Add(v)
	}
	return m
}

func New[T Keyed](cap ...int) KeyedSet[T] {
	return WithCustom(defaultKeyFunc[T], cap[0])
}

func FromSlice[S ~[]T, T Keyed](elems S) KeyedSet[T] {
	return FromWithCustom(defaultKeyFunc[T], elems)
}

// -- Unique methods --

// Identify returns the key for a provided element.
func (ks KeyedSet[T]) Identify(elem T) uint64 {
	return ks.fn(elem)
}

// ContainsKey returns true if there is the key is present in the set.
func (ks KeyedSet[T]) ContainsKey(k uint64) bool {
	_, ok := ks.smap[k]
	return ok
}

// Clone returns a new instance with a shallow copy of the set contents.
func (ks KeyedSet[T]) Clone() KeyedSet[T] {
	return KeyedSet[T]{smap: maps.Clone(ks.smap), fn: ks.fn}
}

// -- SetLike[T] operations --

// Contains returns true if there is an entry matching the element's key.
func (ks KeyedSet[T]) Contains(elem T) bool {
	_, ok := ks.smap[ks.fn(elem)]
	return ok
}

// Len returns the number of elements.
func (ks KeyedSet[T]) Len() int {
	return len(ks.smap)
}

// Range calls yield for each element in the set (implements iter.Seq[T])
func (ks KeyedSet[T]) Range(yield func(T) bool) {
	if yield == nil {
		panic("nil yield function in Set[T].Range")
	}
	for _, v := range ks.smap {
		if !yield(v) {
			return
		}
	}
}

// -- MutableSet[T] operations --

// Add inserts an element into the set
func (ks KeyedSet[T]) Add(elem T) {
	ks.smap[ks.fn(elem)] = elem
}

// Delete removes an element from the set.
func (ks KeyedSet[T]) Delete(elem T) {
	delete(ks.smap, ks.fn(elem))
}
