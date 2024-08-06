package main

import (
	"github.com/stretchr/testify/require"
	"math/rand/v2"
	"testing"
)

func vectorAdditionV0(first, second, dst []uint8) {
	for i := 0; i < len(first); i++ {
		dst[i] = first[i] + second[i]
	}
}

func vectorAdditionV1(first, second, dst []uint8)

func BenchmarkAdd(b *testing.B) {
	b.Run("vectorAdditionV1 (SIMD)", func(b *testing.B) {
		b.StopTimer()
		f, s, dst := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			vectorAdditionV1(f, s, dst)
		}
	})

	b.Run("vectorAdditionV0", func(b *testing.B) {
		b.StopTimer()
		f, s, dst := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			vectorAdditionV0(f, s, dst)
		}
	})
}

func TestVectorAddition(t *testing.T) {
	t.Parallel()

	first := make([]uint8, 160)

	for i := 0; i < 160; i++ {
		first[i] = 10
	}

	second := make([]uint8, 160)

	for i := 0; i < 160; i++ {
		second[i] = 20
	}

	dstV0 := make([]uint8, 160)
	dstV1 := make([]uint8, 160)

	expDst := make([]uint8, 160)

	for i := 0; i < 160; i++ {
		expDst[i] = 30
	}

	vectorAdditionV0(first, second, dstV0)
	vectorAdditionV1(first, second, dstV1)

	require.Equal(t, expDst, dstV1)
	require.Equal(t, dstV1, dstV0)
}

// assume alignment
func getData() ([]uint8, []uint8, []uint8) {
	first := make([]uint8, 1_000_000)

	for i := 0; i < len(first); i++ {
		first[i] = rand.N[uint8](5)
	}

	second := make([]uint8, 1_000_000)

	for i := 0; i < len(second); i++ {
		second[i] = rand.N[uint8](5)
	}

	dst := make([]uint8, 1_000_000)

	return first, second, dst
}
