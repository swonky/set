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

	WithReadLock(func(SetLike[T]))
}

type Cloner[T any] interface {
	Clone() T
}
