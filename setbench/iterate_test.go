package setbench

import "testing"

func BenchmarkRangeFull(b *testing.B) {
	for _, n := range []int{8,32,128,1024,10000} {
		s:=mkSet(n); f:=mkFrozen(n); st:=mkStable(n)
		b.Run("Set/"+itoa(n), func(b *testing.B){
			for i:=0;i<b.N;i++ { c:=0; s.Range(func(v int) bool { c+=v; return true }); sinkInt=c }
		})
		b.Run("Frozen/"+itoa(n), func(b *testing.B){
			for i:=0;i<b.N;i++ { c:=0; f.Range(func(v int) bool { c+=v; return true }); sinkInt=c }
		})
		b.Run("Stable/"+itoa(n), func(b *testing.B){
			for i:=0;i<b.N;i++ { c:=0; st.Range(func(v int) bool { c+=v; return true }); sinkInt=c }
		})
	}
}
