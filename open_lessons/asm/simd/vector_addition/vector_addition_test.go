package main

import (
	"math/rand/v2"
	"testing"
)

func vectorAddition(first, second, dst []uint32)

func vectorAdditionV0(first, second, dst []uint32) {
	for i := 0; i < len(first); i++ {
		dst[i] = first[i] + second[i]
	}
}

func BenchmarkAdd(b *testing.B) {
	b.Run("SIMD vector addition", func(b *testing.B) {
		b.StopTimer()
		f, s, dst := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			vectorAddition(f, s, dst)
		}
	})

	b.Run("simple vector addition", func(b *testing.B) {
		b.StopTimer()
		f, s, dst := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			vectorAdditionV0(f, s, dst)
		}
	})
}

// assume alignment
func getData() ([]uint32, []uint32, []uint32) {
	first := make([]uint32, 1_000_000)

	for i := 0; i < len(first); i++ {
		first[i] = rand.N[uint32](5)
	}

	second := make([]uint32, 1_000_000)

	for i := 0; i < len(second); i++ {
		second[i] = rand.N[uint32](5)
	}

	dst := make([]uint32, 1_000_000)

	return first, second, dst
}
