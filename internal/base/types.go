package base

type SetLike[T any] interface {
	Contains(item T) bool
	Range(func(T) bool)
	Len() int
}

type MutableSet[T any] interface {
	SetLike[T]

	Add(T)
	Delete(T)
}

type ValueSet[T comparable] interface {
	SetLike[T]
}

type LockableSet[T any] interface {
	SetLike[T]

	WithRLock(func(SetLike[T]))
}

type MutableLockableSet[T any] interface {
	MutableSet[T]
	LockableSet[T]

	WithRWLock(func(SetLike[T]))
}

type AsSetter[T comparable] interface {
	Set() Set[T]
}

type Snapshotable[T any] interface {
	Snapshot() SetLike[T]
}
