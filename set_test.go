package set_test

import (
	"testing"

	"github.com/swonky/set"
	"github.com/swonky/set/syncset"
)

func TestCasting(t *testing.T) {
	ss := syncset.New[int](10)
	set.Add(ss, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9)

	if _, ok := any(ss).(set.LockableSet[set.Set[int], int]); !ok {
		t.Fatalf("ss is not a set.LockableSet[set.Set[int], int]")
	}
	if _, ok := any(ss).(set.LockableSet[set.MutableSet[int], int]); !ok {
		t.Fatalf("ss is not a set.LockableSet[set.MutableSet[int], int]")
	}

}

// import (
// 	"iter"
// 	"testing"

// 	"github.com/swonky/set"
// )

// func collect[T comparable](it iter.Seq[T]) set.Set[T] {
// 	out := make(set.Set[T])
// 	it(func(v T) bool {
// 		out[v] = struct{}{}
// 		return true
// 	})
// 	return out
// }

// func seqOf[T any](xs ...T) iter.Seq[T] {
// 	return func(yield func(T) bool) {
// 		for _, x := range xs {
// 			if !yield(x) {
// 				return
// 			}
// 		}
// 	}
// }

// func setEq[T comparable](t *testing.T, got, want set.Set[T]) {
// 	t.Helper()
// 	if got.Len() != want.Len() {
// 		t.Fatalf("len mismatch: got=%d want=%d", got.Len(), want.Len())
// 	}
// 	for v := range want.Range {
// 		if !got.Contains(v) {
// 			t.Fatalf("missing element: %v", v)
// 		}
// 	}
// }

// func asSetLike[T comparable](xs []set.Set[T]) []types.SetLike[T] {
// 	out := make([]types.SetLike[T], len(xs))
// 	for i := range xs {
// 		out[i] = xs[i]
// 	}
// 	return out
// }

// func frozenAsSetLike[T comparable](xs []set.FrozenSet[T]) []types.SetLike[T] {
// 	out := make([]types.SetLike[T], len(xs))
// 	for i := range xs {
// 		out[i] = xs[i]
// 	}
// 	return out
// }

// func syncAsSetLike[T comparable](xs []*set.SyncSet[T]) []types.SetLike[T] {
// 	out := make([]types.SetLike[T], len(xs))
// 	for i := range xs {
// 		out[i] = xs[i]
// 	}
// 	return out
// }

// // func TestNew(t *testing.T) {
// // 	s := set.New(1, 2, 2, 3)

// // 	if s.Len() != 3 {
// // 		t.Fatalf("expected len=3, got %d", s.Len())
// // 	}
// // }

// // func TestFromIter(t *testing.T) {
// // 	src := []int{1, 2, 2, 3}

// // 	s := set.Collect(slices.Values(src))

// // 	if s.Len() != 3 {
// // 		t.Fatalf("expected len=3, got %d", s.Len())
// // 	}
// // }

// // func TestUnion(t *testing.T) {
// // 	a := set.New(1, 2)
// // 	b := set.New(2, 3)

// // 	r := a.Union(b)

// // 	if !r.EqualFunc(set.New(1, 2, 3)) {
// // 		t.Fatalf("unexpected union: %v", r)
// // 	}
// // }

// // func TestUnionInto(t *testing.T) {
// // 	a := set.New(1, 2)
// // 	b := set.New(2, 3)

// // 	a.UnionInto(b)

// // 	if !a.EqualFunc(set.New(1, 2, 3)) {
// // 		t.Fatalf("unexpected result: %v", a)
// // 	}
// // }

// // func TestUnionFn(t *testing.T) {
// // 	a := set.New(1)
// // 	b := set.New(2)
// // 	c := set.New(3)
// // 	sl := []set.SetLike[int]{a, b, c}

// // 	r := set.Unite(sl)

// // 	s := set.Collect(r.Range)

// // 	if !s.EqualFunc(set.New(1, 2, 3)) {
// // 		t.Fatalf("unexpected result: %v", r)
// // 	}
// // }

// // func TestIntersect(t *testing.T) {
// // 	a := set.New(1, 2, 3)
// // 	b := set.New(2, 3, 4)

// // 	r := a.Intersect(b)

// // 	if !r.EqualFunc(set.New(2, 3)) {
// // 		t.Fatalf("unexpected result: %v", r)
// // 	}
// // }

// // func TestIntersectFn(t *testing.T) {
// // 	a := set.New(1, 2, 3)
// // 	b := set.New(2, 3)
// // 	c := set.New(3)

// // 	sl := []set.SetLike[int]{a, b, c}

// // 	r := set.Intersect(sl)

// // 	if !r.AsSet().EqualFunc(set.New(3)) {
// // 		t.Fatalf("unexpected result: %v", r)
// // 	}

// // 	if !r.AsSet().EqualFunc(set.New(3)) {
// // 		t.Fatalf("unexpected result: %v", r)
// // 	}
// // }

// func TestIntersectSort(t *testing.T) {
// 	a := set.New(1, 2, 3, 4)
// 	b := set.New(2, 3)
// 	c := set.New(3)

// 	sl := []set.SetLike[int]{a, b, c}

// 	r := set.Intersect(sl)

// 	for s := range r.Sets() {
// 		t.Logf("%d", s.Len())
// 	}
// }

// // func TestDiff(t *testing.T) {
// // 	a := set.New(1, 2, 3)
// // 	b := set.New(2)

// // 	r := a.Diff(b)

// // 	if !r.EqualFunc(set.New(1, 3)) {
// // 		t.Fatalf("unexpected result: %v", r)
// // 	}
// // }

// // func TestSymDiff(t *testing.T) {
// // 	a := set.New(1, 2)
// // 	b := set.New(2, 3)

// // 	r := a.SymDiff(b)

// // 	if !r.EqualFunc(set.New(1, 3)) {
// // 		t.Fatalf("unexpected result: %v", r)
// // 	}
// // }

// func TestAddCheck(t *testing.T) {
// 	s := set.New[int]()

// 	if found := set.AddCheck(s, 1); found {
// 		t.Fatalf("expected found==false, got true")
// 	}

// 	if found := set.AddCheck(s, 2); found {
// 		t.Fatalf("expected found==false, got true")
// 	}

// 	if found := set.AddCheck(s, 2); !found {
// 		t.Fatalf("expected found==true, got false")
// 	}

// 	if s.Len() != 2 {
// 		t.Fatalf("expected len=2, got %d", s.Len())
// 	}
// }

// func TestAddAll(t *testing.T) {
// 	s := set.New[int]()

// 	set.Append(s, 1, 2, 2)

// 	if s.Len() != 2 {
// 		t.Fatalf("expected len=2, got %d", s.Len())
// 	}
// }

// func TestHas_Delete(t *testing.T) {
// 	s := set.New(1, 2)

// 	if !s.Contains(1) {
// 		t.Fatalf("expected to have 1")
// 	}

// 	s.Delete(1)

// 	if s.Contains(1) {
// 		t.Fatalf("expected 1 to be deleted")
// 	}
// }

// func TestPop(t *testing.T) {
// 	s := set.New(1)

// 	v, ok := set.Pop(s)

// 	if !ok {
// 		t.Fatalf("expected ok")
// 	}

// 	if v != 1 {
// 		t.Fatalf("unexpected value: %v", v)
// 	}

// 	if !s.IsEmpty() {
// 		t.Fatalf("expected empty after pop")
// 	}

// 	_, ok = set.Pop(s)
// 	if ok {
// 		t.Fatalf("expected empty pop to be false")
// 	}
// }

// func TestClear(t *testing.T) {
// 	s := set.New(1, 2, 3)

// 	set.Clear(s)

// 	if !s.IsEmpty() {
// 		t.Fatalf("expected empty")
// 	}
// }

// func TestClone(t *testing.T) {
// 	s := set.New(1, 2)

// 	c := s.Clone()
// 	c.Add(3)

// 	if s.Contains(3) {
// 		t.Fatalf("clone should not affect original")
// 	}
// }

// func TestLen_IsEmpty(t *testing.T) {
// 	s := set.New[int]()

// 	if !s.IsEmpty() {
// 		t.Fatalf("expected empty")
// 	}

// 	s.Add(1)

// 	if s.IsEmpty() || s.Len() != 1 {
// 		t.Fatalf("unexpected state")
// 	}
// }

// func TestSubsetSuperset(t *testing.T) {
// 	a := set.New(1, 2)
// 	b := set.New(1, 2, 3)

// 	if !a.IsSubsetOf(b) {
// 		t.Fatalf("expected subset")
// 	}

// 	if !b.IsSupersetOf(a) {
// 		t.Fatalf("expected superset")
// 	}
// }

// func TestEqualFunc(t *testing.T) {
// 	a := set.New(1, 2)
// 	b := set.New(2, 1)

// 	if !a.EqualFunc(b) {
// 		t.Fatalf("expected equal")
// 	}
// }

// func TestIter_AsSlice(t *testing.T) {
// 	s := set.New(1, 2, 3)

// 	out := s.AsSlice()

// 	if len(out) != 3 {
// 		t.Fatalf("expected len=3")
// 	}

// 	for _, v := range out {
// 		if !s.Contains(v) {
// 			t.Fatalf("unexpected value: %v", v)
// 		}
// 	}
// }

// func TestString(t *testing.T) {
// 	s := set.New(1, 2)

// 	str := s.String()

// 	if len(str) == 0 || str[0] != '{' {
// 		t.Fatalf("unexpected string: %s", str)
// 	}
// }

// func TestFilter(t *testing.T) {
// 	s := set.New(1, 2, 3, 4)

// 	r := s.Filter(func(x int) bool { return x%2 == 0 })

// 	if !r.EqualFunc(set.New(2, 4)) {
// 		t.Fatalf("unexpected result: %v", r)
// 	}
// }

// func TestAny_All(t *testing.T) {
// 	s := set.New(2, 4)

// 	if !s.AllFunc(func(x int) bool { return x%2 == 0 }) {
// 		t.Fatalf("expected all even")
// 	}

// 	if !s.AnyFunc(func(x int) bool { return x == 2 }) {
// 		t.Fatalf("expected any match")
// 	}
// }

// func TestFind(t *testing.T) {
// 	s := set.New(1, 2, 3)

// 	v, ok := s.Find(func(x int) bool { return x == 2 })

// 	if !ok || v != 2 {
// 		t.Fatalf("unexpected result: %v %v", v, ok)
// 	}

// 	_, ok = s.Find(func(x int) bool { return x == 99 })
// 	if ok {
// 		t.Fatalf("expected not found")
// 	}
// }

// func TestFirst(t *testing.T) {
// 	s := set.New(1)

// 	v, ok := s.First()

// 	if !ok || v != 1 {
// 		t.Fatalf("unexpected result")
// 	}

// 	set.Clear(s)

// 	_, ok = s.First()
// 	if ok {
// 		t.Fatalf("expected empty")
// 	}
// }

// func TestCollect(t *testing.T) {
// 	s := set.Collect(seqOf(1, 2, 2, 3))
// 	setEq(t, s, set.New(1, 2, 3))
// }

// func TestAccumulate(t *testing.T) {
// 	a := set.New(1, 2)
// 	b := set.New(2, 3)
// 	c := set.New(3, 4)

// 	var results []set.Set[int]
// 	for v := range set.Accumulate(
// 		[]set.Set[int]{a, b, c},
// 		func(x, y set.Set[int]) set.Set[int] {
// 			return x.Union(y)
// 		},
// 	) {
// 		results = append(results, v)
// 	}

// 	if len(results) != 2 {
// 		t.Fatalf("expected 2 results, got %d: %v", len(results), results)
// 	}

// 	setEq(t, results[0], set.New(1, 2, 3))
// 	setEq(t, results[1], set.New(1, 2, 3, 4))
// }

// func TestAccumulateWhileStops(t *testing.T) {
// 	a := set.New(1)
// 	b := set.New(2)
// 	c := set.New(3)

// 	count := 0
// 	for range set.AccumulateTry([]set.Set[int]{a, b, c},
// 		func(x, y set.Set[int]) (set.Set[int], bool) {
// 			count++
// 			return x.Union(y), false
// 		},
// 	) {
// 	}

// 	if count != 1 {
// 		t.Fatalf("expected early stop after 1 step, got %d", count)
// 	}
// }

// func TestReduce(t *testing.T) {
// 	a := set.New(1)
// 	b := set.New(2)
// 	c := set.New(3)

// 	got := set.Reduce([]set.Set[int]{a, b, c},
// 		func(x, y set.Set[int]) set.Set[int] {
// 			return x.Union(y)
// 		},
// 	)

// 	setEq(t, got, set.New(1, 2, 3))
// }

// func TestReduceEmpty(t *testing.T) {
// 	got := set.Reduce[int](nil,
// 		func(x, y set.Set[int]) set.Set[int] { return x },
// 	)
// 	if !got.IsEmpty() {
// 		t.Fatalf("expected empty set")
// 	}
// }

// func TestReduceTryStops(t *testing.T) {
// 	a := set.New(1)
// 	b := set.New(2)
// 	c := set.New(3)

// 	count := 0
// 	got := set.ReduceTry([]set.Set[int]{a, b, c},
// 		func(x, y set.Set[int]) (set.Set[int], bool) {
// 			count++
// 			return x.Union(y), count < 2
// 		},
// 	)

// 	setEq(t, got, set.New(1, 2, 3))
// 	if count != 2 {
// 		t.Fatalf("expected 2 steps, got %d", count)
// 	}
// }

// func TestReduceWhileStops(t *testing.T) {
// 	a := set.New(1)
// 	b := set.New(2)
// 	c := set.New(3)

// 	got := set.ReduceWhile(
// 		[]set.Set[int]{a, b, c},
// 		func(x, y set.Set[int]) set.Set[int] { return x.Union(y) },
// 		func(s set.Set[int]) bool { return s.Len() < 3 },
// 	)

// 	setEq(t, got, set.New(1, 2, 3))
// }

// func TestIntersectEarlyEmpty(t *testing.T) {

// 	sl := []set.SetLike[int]{set.New(1), set.New(2), set.New(3)}
// 	got := set.Intersect(sl)
// 	if got.Len() != 0 {
// 		t.Fatalf("expected empty intersection")
// 	}
// }

// // func TestUnionIter(t *testing.T) {
// // 	a := set.New(1, 2)
// // 	b := set.New(2, 3)
// // 	c := set.New(3, 4)

// // 	var out []int
// // 	for v := range set.UnionIter(a, b, c) {
// // 		out = append(out, v)
// // 	}

// // 	setEq(t, set.New(out...), set.New(1, 2, 3, 4))
// // }

// // }

// // func TestUnionIter2(t *testing.T) {
// // 	a := set.New(1, 2)
// // 	c := set.New(2, 3)
// // 	d := set.New(3, 4)
// // 	sl := []set.SetLike[int]{a, c, d}

// // 	var out []int
// // 	for v := range set.UnionIter2(sl) {
// // 		out = append(out, v)
// // 	}

// // 	setEq(t, set.New(out...), set.New(1, 2, 3, 4))
// // }

// // func TestUnionIterOrder(t *testing.T) {
// // 	a := set.New(1, 2)
// // 	b := set.New(2, 3)

// // 	var out []int
// // 	for v := range set.UnionIter(a, b) {
// // 		out = append(out, v)
// // 	}

// // 	expected := []int{1, 2, 3}
// // 	if len(out) != len(expected) {
// // 		t.Fatalf("length mismatch: got=%v want=%v", out, expected)
// // 	}
// // 	for i := range expected {
// // 		if out[i] != expected[i] {
// // 			t.Fatalf("order mismatch: got=%v want=%v", out, expected)
// // 		}
// // 	}
// // }
