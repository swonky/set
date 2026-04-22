package set

import (
	"github.com/swonky/set/internal/base"
	"github.com/swonky/set/lazyset"
)

type SetLike[T any] = base.SetLike[T]
type MutableSet[T any] = base.MutableSet[T]
type LockableSet[T any] = base.LockableSet[T]
type ValueSet[T comparable] = base.ValueSet[T]
type AsSetter[T comparable] = base.AsSetter[T]

type Set[T comparable] = base.Set[T]

type Intersection[S SetLike[T], T any] = lazyset.Intersection[S, T]
type Union[S SetLike[T], T any] = lazyset.Union[S, T]
