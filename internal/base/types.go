package base

type SetLike[T any] interface {
	Contains(item T) bool
	Range(func(T) bool)
	Len() int
}

type ValueSet[T comparable] interface {
	SetLike[T]
}

type MutableSet[T any] interface {
	SetLike[T]

	Add(T)
	Delete(T)
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
	AsSet() Set[T]
}

type Snapshotable[T any] interface {
	Snapshot() SetLike[T]
}
