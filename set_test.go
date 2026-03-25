package set_test

import (
	"slices"
	"testing"

	"github.com/swonky/set"
)

func TestNew(t *testing.T) {
	s := set.New(1, 2, 2, 3)

	if s.Len() != 3 {
		t.Fatalf("expected len=3, got %d", s.Len())
	}
}

func TestFromIter(t *testing.T) {
	src := []int{1, 2, 2, 3}

	s := set.FromIter(slices.Values(src))

	if s.Len() != 3 {
		t.Fatalf("expected len=3, got %d", s.Len())
	}
}

func TestUnion(t *testing.T) {
	a := set.New(1, 2)
	b := set.New(2, 3)

	r := a.Union(b)

	if !r.Equal(set.New(1, 2, 3)) {
		t.Fatalf("unexpected union: %v", r)
	}
}

func TestUnionInto(t *testing.T) {
	a := set.New(1, 2)
	b := set.New(2, 3)

	a.UnionInto(b)

	if !a.Equal(set.New(1, 2, 3)) {
		t.Fatalf("unexpected result: %v", a)
	}
}

func TestUnionAll(t *testing.T) {
	a := set.New(1)
	b := set.New(2)
	c := set.New(3)

	r := set.UnionAll(a, b, c)

	if !r.Equal(set.New(1, 2, 3)) {
		t.Fatalf("unexpected result: %v", r)
	}
}

func TestIntersect(t *testing.T) {
	a := set.New(1, 2, 3)
	b := set.New(2, 3, 4)

	r := a.Intersect(b)

	if !r.Equal(set.New(2, 3)) {
		t.Fatalf("unexpected result: %v", r)
	}
}

func TestIntersectAll(t *testing.T) {
	a := set.New(1, 2, 3)
	b := set.New(2, 3)
	c := set.New(3)

	r := set.IntersectAll(a, b, c)

	if !r.Equal(set.New(3)) {
		t.Fatalf("unexpected result: %v", r)
	}
}

func TestDiff(t *testing.T) {
	a := set.New(1, 2, 3)
	b := set.New(2)

	r := a.Diff(b)

	if !r.Equal(set.New(1, 3)) {
		t.Fatalf("unexpected result: %v", r)
	}
}

func TestSymmetricDiff(t *testing.T) {
	a := set.New(1, 2)
	b := set.New(2, 3)

	r := a.SymmetricDiff(b)

	if !r.Equal(set.New(1, 3)) {
		t.Fatalf("unexpected result: %v", r)
	}
}

func TestAdd(t *testing.T) {
	s := set.New[int]()

	s.Add(1, 2, 2)

	if s.Len() != 2 {
		t.Fatalf("expected len=2, got %d", s.Len())
	}
}

func TestHas_Delete(t *testing.T) {
	s := set.New(1, 2)

	if !s.Has(1) {
		t.Fatalf("expected to have 1")
	}

	s.Delete(1)

	if s.Has(1) {
		t.Fatalf("expected 1 to be deleted")
	}
}

func TestPop(t *testing.T) {
	s := set.New(1)

	v, ok := s.Pop()

	if !ok {
		t.Fatalf("expected ok")
	}

	if v != 1 {
		t.Fatalf("unexpected value: %v", v)
	}

	if !s.IsEmpty() {
		t.Fatalf("expected empty after pop")
	}

	_, ok = s.Pop()
	if ok {
		t.Fatalf("expected empty pop to be false")
	}
}

func TestClear(t *testing.T) {
	s := set.New(1, 2, 3)

	s.Clear()

	if !s.IsEmpty() {
		t.Fatalf("expected empty")
	}
}

func TestClone(t *testing.T) {
	s := set.New(1, 2)

	c := s.Clone()
	c.Add(3)

	if s.Has(3) {
		t.Fatalf("clone should not affect original")
	}
}

func TestLen_IsEmpty(t *testing.T) {
	s := set.New[int]()

	if !s.IsEmpty() {
		t.Fatalf("expected empty")
	}

	s.Add(1)

	if s.IsEmpty() || s.Len() != 1 {
		t.Fatalf("unexpected state")
	}
}

func TestSubsetSuperset(t *testing.T) {
	a := set.New(1, 2)
	b := set.New(1, 2, 3)

	if !a.IsSubsetOf(b) {
		t.Fatalf("expected subset")
	}

	if !b.IsSupersetOf(a) {
		t.Fatalf("expected superset")
	}
}

func TestEqual(t *testing.T) {
	a := set.New(1, 2)
	b := set.New(2, 1)

	if !a.Equal(b) {
		t.Fatalf("expected equal")
	}
}

func TestIter_AsSlice(t *testing.T) {
	s := set.New(1, 2, 3)

	out := s.AsSlice()

	if len(out) != 3 {
		t.Fatalf("expected len=3")
	}

	for _, v := range out {
		if !s.Has(v) {
			t.Fatalf("unexpected value: %v", v)
		}
	}
}

func TestString(t *testing.T) {
	s := set.New(1, 2)

	str := s.String()

	if len(str) == 0 || str[0] != '{' {
		t.Fatalf("unexpected string: %s", str)
	}
}

func TestFilter(t *testing.T) {
	s := set.New(1, 2, 3, 4)

	r := s.Filter(func(x int) bool { return x%2 == 0 })

	if !r.Equal(set.New(2, 4)) {
		t.Fatalf("unexpected result: %v", r)
	}
}

func TestAny_All(t *testing.T) {
	s := set.New(2, 4)

	if !s.All(func(x int) bool { return x%2 == 0 }) {
		t.Fatalf("expected all even")
	}

	if !s.Any(func(x int) bool { return x == 2 }) {
		t.Fatalf("expected any match")
	}
}

func TestFind(t *testing.T) {
	s := set.New(1, 2, 3)

	v, ok := s.Find(func(x int) bool { return x == 2 })

	if !ok || v != 2 {
		t.Fatalf("unexpected result: %v %v", v, ok)
	}

	_, ok = s.Find(func(x int) bool { return x == 99 })
	if ok {
		t.Fatalf("expected not found")
	}
}

func TestFirst(t *testing.T) {
	s := set.New(1)

	v, ok := s.First()

	if !ok || v != 1 {
		t.Fatalf("unexpected result")
	}

	s.Clear()

	_, ok = s.First()
	if ok {
		t.Fatalf("expected empty")
	}
}
