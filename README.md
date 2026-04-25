# set

Package `set` provides generic set interfaces, specialized implementations, and reusable operations for collections of unique values.

[![Go Reference](https://pkg.go.dev/badge/github.com/swonky/set.svg)](https://pkg.go.dev/github.com/swonky/set)

A set is easy to write until it needs to be good.

Writing a set is simple. Writing one that behaves well around the compiler, allocator, garbage collector, iteration costs, concurrency, and real workloads is less simple. The point of `set` to provide implementations for common uses ergonomic, consistent, and cheap.

## Features

## Features

- Generic set implementations:
  - `set.Set[T]` — general-purpose mutable hash set.
  - `frozenset.FrozenSet[T]` — immutable read-optimized set.
  - `syncset.SyncSet[T]` — concurrency-safe mutable set.
  - `stableset.StableSet[T]` — insertion-ordered set with stable iteration.
  - `keyedset.KeyedSet[T]` — keyed set for non-comparable values.
- Optimised typed set implementations:
  - `bitset.BitSet` - mutable set of non-negative integers.
- Common interfaces (`set.SetLike[T]`, `set.MutableSet[T]`) for interoperability across implementations.
- Standard-library-style set operations including `Diff`/`SymmDiff`, `Union`, `Intersect`, `Filter`, `Transform`, and [more...]()
- Lazy set views for union and intersection to avoid unnecessary materialization.
- Generic constructor and convertor functions: `New`, `FromSlice`, `FromSetLike`, `Collect`.
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