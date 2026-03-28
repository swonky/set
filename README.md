# set

Package set provides a generic `Set` type for working with unordered collections of unique, comparable values.

[![Go Reference](https://pkg.go.dev/badge/github.com/swonky/set.svg)](https://pkg.go.dev/github.com/swonky/set)

## Features

- Three set types: `Set[T]` (mutable), `FrozenSet[T]` (immutable), `SyncSet[T]` (thread-safe mutable).
- Supports lazy iteration using Go 1.22+ `iter` package.
- Composable set operations using reducers and accumulation patterns (`set.Reduce`, `set.Accumulate`).
- Shared `set.SetLike` interface enabling interoperability across set types.

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

if s.Has(2) {
    // ...
}
```
