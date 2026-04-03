package set

import (
	"github.com/swonky/set/frozenset"
	"github.com/swonky/set/internal/base"
	"github.com/swonky/set/keyedset"
	"github.com/swonky/set/syncset"
)

type SetLike[T any] = base.SetLike[T]
type MutableSet[T any] = base.MutableSet[T]
type LockableSet[T any] = base.LockableSet[T]
type ValueSet[T comparable] = base.ValueSet[T]

type Cloner[T any] = base.Cloner[T]

type Set[T comparable] = base.Set[T]
type FrozenSet[T comparable] = frozenset.FrozenSet[T]
type SyncSet[T comparable] = syncset.SyncSet[T]

type KeyedSet[T any] = keyedset.KeyedSet[T]
type Keyed = keyedset.Keyed
type Key = keyedset.Key
