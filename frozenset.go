package set

import "iter"

var _ SetLike[int] = FrozenSet[int]{}

type FrozenSet[T comparable] struct {
	s Set[T]
}

func (fs FrozenSet[T]) All(fn func(T) bool) bool                 { return fs.s.All(fn) }
func (fs FrozenSet[T]) Any(fn func(T) bool) bool                 { return fs.s.Any(fn) }
func (fs FrozenSet[T]) AsSet() Set[T]                            { return fs.s.Clone() }
func (fs FrozenSet[T]) AsSlice() []T                             { return fs.s.AsSlice() }
func (fs FrozenSet[T]) Clone() FrozenSet[T]                      { return fs.s.Clone().Freeze() }
func (fs FrozenSet[T]) Diff(o FrozenSet[T]) FrozenSet[T]         { return fs.s.Diff(o.s).Freeze() }
func (fs FrozenSet[T]) Equal(o FrozenSet[T]) bool                { return fs.s.Equal(o.s) }
func (fs FrozenSet[T]) Filter(fn func(T) bool) FrozenSet[T]      { return fs.s.Filter(fn).Freeze() }
func (fs FrozenSet[T]) Find(fn func(T) bool) (T, bool)           { return fs.s.Find(fn) }
func (fs FrozenSet[T]) First() (T, bool)                         { return fs.s.First() }
func (fs FrozenSet[T]) Has(item T) bool                          { return fs.s.Has(item) }
func (fs FrozenSet[T]) HasAll(item ...T) bool                    { return fs.s.HasAll(item...) }
func (fs FrozenSet[T]) HasAny(item ...T) bool                    { return fs.s.HasAny(item...) }
func (fs FrozenSet[T]) Intersect(o FrozenSet[T]) FrozenSet[T]    { return fs.s.Intersect(o.s).Freeze() }
func (fs FrozenSet[T]) IntersectIter(o FrozenSet[T]) iter.Seq[T] { return fs.s.IntersectIter(o.s) }
func (fs FrozenSet[T]) IsEmpty() bool                            { return fs.s.IsEmpty() }
func (fs FrozenSet[T]) IsSubsetOf(o FrozenSet[T]) bool           { return fs.s.IsSubsetOf(o.s) }
func (fs FrozenSet[T]) IsSupersetOf(o FrozenSet[T]) bool         { return fs.s.IsSupersetOf(o.s) }
func (fs FrozenSet[T]) Iter() iter.Seq[T]                        { return fs.s.Iter() }
func (fs FrozenSet[T]) Len() int                                 { return fs.s.Len() }
func (fs FrozenSet[T]) String() string                           { return fs.s.String() }
func (fs FrozenSet[T]) SymDiff(o FrozenSet[T]) FrozenSet[T]      { return fs.s.SymDiff(o.s).Freeze() }
func (fs FrozenSet[T]) Union(o FrozenSet[T]) FrozenSet[T]        { return fs.s.Union(o.s).Freeze() }
func (fs FrozenSet[T]) UnionIter(o FrozenSet[T]) iter.Seq[T]     { return fs.s.UnionIter(o.s) }
