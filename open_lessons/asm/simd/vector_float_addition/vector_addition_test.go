package main

import (
	"github.com/stretchr/testify/require"
	"math"
	"math/rand/v2"
	"testing"
)

func vectorFloatAdditionV0(first, second, dst []float32) {
	for i := 0; i < len(first); i++ {
		dst[i] = first[i] + second[i]
	}
}

func vectorFloatAdditionV1(first, second, dst []float32)

func BenchmarkFloatAdd(b *testing.B) {
	b.Run("vectorFloatAdditionV1 (SIMD)", func(b *testing.B) {
		b.StopTimer()
		f, s, dst := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			vectorFloatAdditionV1(f, s, dst)
		}
	})

	b.Run("vectorFloatAdditionV0", func(b *testing.B) {
		b.StopTimer()
		f, s, dst := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			vectorFloatAdditionV0(f, s, dst)
		}
	})
}

func TestVectorAddition(t *testing.T) {
	t.Parallel()

	first := make([]float32, 160)

	for i := 0; i < 160; i++ {
		first[i] = 10.5
	}

	second := make([]float32, 160)

	for i := 0; i < 160; i++ {
		second[i] = 20.4
	}

	dstV0 := make([]float32, 160)
	dstV1 := make([]float32, 160)

	expDst := make([]float32, 160)

	for i := 0; i < 160; i++ {
		expDst[i] = 30.9
	}

	vectorFloatAdditionV0(first, second, dstV0)
	vectorFloatAdditionV1(first, second, dstV1)

	for i := 0; i < 160; i++ {
		require.True(t, math.Abs(float64(expDst[i]-dstV1[i])) < 0.000001)
		require.True(t, math.Abs(float64(expDst[i]-dstV0[i])) < 0.000001)
	}
}

// assume alignment
func getData() ([]float32, []float32, []float32) {
	first := make([]float32, 1_000_000)

	for i := 0; i < len(first); i++ {
		first[i] = float32(rand.N[int](5)) + 0.5
	}

	second := make([]float32, 1_000_000)

	for i := 0; i < len(second); i++ {
		second[i] = float32(rand.N[int](5)) + 0.5
	}

	dst := make([]float32, 1_000_000)

	return first, second, dst
}
