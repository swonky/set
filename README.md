# set

Package set provides a generic `Set` type for working with unordered collections of unique, comparable values.

[![Go Reference](https://pkg.go.dev/badge/github.com/swonky/set.svg)](https://pkg.go.dev/github.com/swonky/set)

## Features

- Supports both non-mutating (returning new sets) and mutating (in-place) operations.
- Includes helpers for querying elements (eg. `Any`, `All`, `Find`).
- Includes functional helpers (eg. `Reduce`, `ReduceWhile`, `ReduceTry`)
- Supports iterator patterns via `iter.Seq` without allocating intermediate slices (eg. `Collect`, `Chain`, `Accumulate`, `AccumulateTry`).

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
