package setbench

import "testing"

func BenchmarkContainsHit(b *testing.B) {
	for _, n := range []int{8,32,128,1024,10000} {
		s:=mkSet(n); f:=mkFrozen(n); st:=mkStable(n)
		target:=n-1
		b.Run("Set/"+itoa(n), func(b *testing.B){ for i:=0;i<b.N;i++ { sinkBool=s.Contains(target) }})
		b.Run("Frozen/"+itoa(n), func(b *testing.B){ for i:=0;i<b.N;i++ { sinkBool=f.Contains(target) }})
		b.Run("Stable/"+itoa(n), func(b *testing.B){ for i:=0;i<b.N;i++ { sinkBool=st.Contains(target) }})
	}
}

func BenchmarkContainsMiss(b *testing.B) {
	for _, n := range []int{8,32,128,1024,10000} {
		s:=mkSet(n); f:=mkFrozen(n); st:=mkStable(n)
		target:=-1
		b.Run("Set/"+itoa(n), func(b *testing.B){ for i:=0;i<b.N;i++ { sinkBool=s.Contains(target) }})
		b.Run("Frozen/"+itoa(n), func(b *testing.B){ for i:=0;i<b.N;i++ { sinkBool=f.Contains(target) }})
		b.Run("Stable/"+itoa(n), func(b *testing.B){ for i:=0;i<b.N;i++ { sinkBool=st.Contains(target) }})
	}
}
