package setbench

import (
	"runtime"
	"sync/atomic"
	"testing"

	set "github.com/swonky/set"
	syncset "github.com/swonky/set/syncset"
)

func mkSync(n int) *syncset.SyncSet[set.Set[int], int] {
	s := syncset.New[int](n)
	for i := 0; i < n; i++ {
		s.Add(i)
	}
	return s
}

func BenchmarkSyncSetContains(b *testing.B) {
	for _, n := range []int{128, 1024, 10000} {
		s := mkSync(n)
		b.Run(itoa(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sinkBool = s.Contains(n - 1)
			}
		})
	}
}

func BenchmarkSyncSetContainsParallel(b *testing.B) {
	for _, n := range []int{128, 1024, 10000} {
		s := mkSync(n)
		b.Run(itoa(n), func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					sinkBool = s.Contains(n - 1)
				}
			})
		})
	}
}

func benchMixed(b *testing.B, readsPerWrite int, n int) {
	s := mkSync(n)
	var ctr atomic.Uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x := ctr.Add(1)
			if x%uint64(readsPerWrite+1) != 0 {
				sinkBool = s.Contains(int(x % uint64(n)))
			} else {
				v := int(x % uint64(n))
				s.Add(v + n)
				s.Delete(v + n)
			}
		}
	})
}

func BenchmarkSyncSetMixed95_5(b *testing.B) {
	for _, n := range []int{1024, 10000} {
		b.Run(itoa(n), func(b *testing.B) { benchMixed(b, 19, n) })
	}
}
func BenchmarkSyncSetMixed80_20(b *testing.B) {
	for _, n := range []int{1024, 10000} {
		b.Run(itoa(n), func(b *testing.B) { benchMixed(b, 4, n) })
	}
}
func BenchmarkSyncSetMixed50_50(b *testing.B) {
	for _, n := range []int{1024, 10000} {
		b.Run(itoa(n), func(b *testing.B) { benchMixed(b, 1, n) })
	}
}

func BenchmarkSyncSetRange(b *testing.B) {
	for _, n := range []int{128, 1024, 10000} {
		s := mkSync(n)
		b.Run(itoa(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s.Range(func(v int) bool { return true })
			}
		})
	}
}

func BenchmarkSyncSetRangeWithReaders(b *testing.B) {
	for _, n := range []int{1024, 10000} {
		s := mkSync(n)
		b.Run(itoa(n), func(b *testing.B) {
			done := make(chan struct{})
			for i := 0; i < runtime.NumCPU(); i++ {
				go func() {
					for {
						select {
						case <-done:
							return
						default:
							s.Contains(n - 1)
						}
					}
				}()
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Range(func(v int) bool { return true })
			}
			close(done)
		})
	}
}
