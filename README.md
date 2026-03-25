# set

A small, idiomatic generic set implementation for Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/swonky/set.svg)](https://pkg.go.dev/github.com/swonky/set)

## Features

- Provides `Set`, a collection of unique comparable values.
- Uses map[T]struct{} internally for constant-time add, remove, and membership checks.
- Supports both non-mutating operations (returning new sets) and in-place variants.
- Includes helpers for querying elements (e.g. Any, All, Find).
- Supports on-demand iteration via `iter.Seq` without allocating intermediate slices.

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
