package setbench

import (
	set "github.com/swonky/set"
	"github.com/swonky/set/frozenset"
	"github.com/swonky/set/keyedset"
	"github.com/swonky/set/stableset"
)

type benchItem struct {
	id uint64
	v  string
}

func (b benchItem) Key() uint64 { return b.id }

var sinkBool bool
var sinkInt int
var sinkAny any

func ints(n int) []int {
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = i
	}
	return out
}
func intsDup(n int) []int {
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = i % 32
	}
	return out
}
func items(n int) []benchItem {
	out := make([]benchItem, n)
	for i := 0; i < n; i++ {
		out[i] = benchItem{id: uint64(i), v: "x"}
	}
	return out
}

func mkSet(n int) set.Set[int]                   { return set.FromSlice(ints(n)) }
func mkFrozen(n int) frozenset.FrozenSet[int]    { return frozenset.FromSlice(ints(n)) }
func mkStable(n int) *stableset.StableSet[int]   { return stableset.FromSlice(ints(n)) }
func mkIndex(n int) keyedset.KeyedSet[benchItem] { return keyedset.FromSlice(items(n)) }
