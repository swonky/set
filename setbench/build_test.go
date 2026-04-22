package setbench

import (
	"testing"
	set "github.com/swonky/set"
	"github.com/swonky/set/frozenset"
	"github.com/swonky/set/stableset"
)

func BenchmarkBuildFromSlice(b *testing.B) {
	for _, n := range []int{0,8,32,128,1024,10000} {
		data := ints(n)
		b.Run("Set/"+itoa(n), func(b *testing.B){
			b.ReportAllocs()
			for i:=0;i<b.N;i++ { sinkAny = set.FromSlice(data) }
		})
		b.Run("Frozen/"+itoa(n), func(b *testing.B){
			b.ReportAllocs()
			for i:=0;i<b.N;i++ { sinkAny = frozenset.FromSlice(data) }
		})
		b.Run("Stable/"+itoa(n), func(b *testing.B){
			b.ReportAllocs()
			for i:=0;i<b.N;i++ { sinkAny = stableset.FromSlice(data) }
		})
	}
}
