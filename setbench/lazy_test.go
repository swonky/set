package setbench

// func BenchmarkUnionContains(b *testing.B) {
// 	for _, n := range []int{32, 128, 1024, 10000} {
// 		u1 := set.Unite(mkSet(n), mkFrozen(n))
// 		u2 := set.Unite(mkStable(n), mkSet(n))
// 		b.Run("SetFrozen/"+itoa(n), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				sinkBool = u1.Contains(n - 1)
// 			}
// 		})
// 		b.Run("StableSet/"+itoa(n), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				sinkBool = u2.Contains(n - 1)
// 			}
// 		})
// 	}
// }

// func BenchmarkIntersectionContains(b *testing.B) {
// 	for _, n := range []int{32, 128, 1024, 10000} {
// 		x1 := set.Intersect(mkSet(n), mkFrozen(n))
// 		b.Run("SetFrozen/"+itoa(n), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				sinkBool = x1.Contains(n - 1)
// 			}
// 		})
// 	}
// }
