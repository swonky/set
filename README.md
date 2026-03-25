# set

A small, idiomatic generic Set for Go.

## Features

- `map[T]struct{}`-backed (O(1) operations)
- Pure + mutating set operations (`Union`, `UnionInto`, etc.)
- Predicate helpers (`Any`, `All`, `Find`)
- Iterator support (`iter.Seq`)

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
