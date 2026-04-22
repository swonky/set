# set

Package `set` provides generic set interfaces, specialized implementations, and reusable operations for collections of unique values.

[![Go Reference](https://pkg.go.dev/badge/github.com/swonky/set.svg)](https://pkg.go.dev/github.com/swonky/set)

A set is easy to write until it needs to be good.

Writing a set is simple. Writing one that behaves well around the compiler, allocator, garbage collector, iteration costs, concurrency, and real workloads is less simple. The point of `set` to provide implementations for common uses ergonomic, consistent, and cheap.

## Features

## Features

- Multiple generic set implementations with shared APIs:
  - `set.Set[T]` — general-purpose mutable hash set
  - `frozenset.FrozenSet[T]` — immutable read-optimized set
  - `syncset.SyncSet[T]` — concurrency-safe mutable set
  - `stableset.StableSet[T]` — insertion-ordered set with stable iteration
  - `keyedset.KeyedSet[T]` — keyed set for non-comparable values
- Common interfaces (`set.SetLike[T]`, `set.MutableSet[T]`) for interoperability across implementations.
- Standard-library-style set operations including diff, union, intersection, symmetric difference, filtering, transforms, equality checks, predicates, and selection helpers.
- Lazy set views for union and intersection to avoid unnecessary materialization.
- Generic constructors and conversions: `New`, `FromSlice`, `FromSetLike`, `Collect`.
- Allocation-conscious implementations designed for practical workloads and hot paths.
- Go 1.22+ iterator support via `iter.Seq`.

## Install

```bash
go get github.com/swonky/set
```

## Usage

```go
s := set.New(1, 2, 3)
s.Add(4)

t := set.New(3, 4, 5)

u := s.Union(t)     // {1,2,3,4,5}
s.UnionInto(t)      // mutate s

if s.Contains(2) {
    // ...
}
```
## Types

### Set[T]
### KeyedSet[T]