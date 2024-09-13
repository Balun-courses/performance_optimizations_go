package main

import (
	"testing"
)

func BenchmarkAppend(b *testing.B) {
	b.Run("index", func(b *testing.B) {
		b.ReportAllocs()

		for j := 0; j < b.N; j++ {
			b.StopTimer()
			a := make([]int64, 100_000)
			b.StartTimer()

			for i := 0; i < len(a); i++ {
				a[i] = int64(i)
			}
		}
	})

	a := make([]int64, 1024<<20)

	b :=

		b.Run("append", func(b *testing.B) {
			b.ReportAllocs()

			for j := 0; j < b.N; j++ {
				b.StopTimer()
				a := make([]int64, 0, 100_000)
				b.StartTimer()

				for i := 0; i < cap(a); i++ {
					a = append(a, int64(i))
				}
			}
		})
}
