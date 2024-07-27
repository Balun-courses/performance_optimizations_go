package main

import (
	"asm/simd/slice_contains/cgo"
	"github.com/stretchr/testify/require"
	"math/rand/v2"
	"slices"
	"testing"
)

func SliceContainsV0(s []uint8, target uint8) bool {
	return slices.Contains(s, target)
}

func SliceContainsV1(s []uint8, target uint8) bool

func BenchmarkSliceContains(b *testing.B) {
	b.Run("SliceContainsV1 (SIMD)", func(b *testing.B) {
		b.StopTimer()
		s, target := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			SliceContainsV1(s, target)
		}

	})

	b.Run("SliceContainsV0", func(b *testing.B) {
		b.StopTimer()
		s, target := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			SliceContainsV0(s, target)
		}
	})

	b.Run("SliceContainsCgo", func(b *testing.B) {
		b.StopTimer()
		s, target := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			cgo.SliceContains(s, target)
		}
	})
}

func TestSliceContains(t *testing.T) {
	t.Parallel()

	t.Run("latest match", func(t *testing.T) {
		t.Parallel()
		s, target := getData()
		require.Equal(t, target, s[len(s)-1])

		// always last match
		require.False(t, slices.Contains(s[:len(s)-1], target))

		simd := SliceContainsV0(s, target)
		simple := SliceContainsV1(s, target)
		cgoResult := cgo.SliceContains(s, target)

		require.True(t, simd)
		require.True(t, simple)
		require.True(t, cgoResult)
	})

	t.Run("middle match", func(t *testing.T) {
		t.Parallel()

		s, target := getData()
		s = s[:len(s)-1]
		s[len(s)/2] = target

		simd := SliceContainsV0(s, target)
		simple := SliceContainsV1(s, target)
		cgoResult := cgo.SliceContains(s, target)

		require.True(t, simd)
		require.True(t, simple)
		require.True(t, cgoResult)
	})

	t.Run("no match", func(t *testing.T) {
		t.Parallel()

		s, target := getData()
		s[len(s)-1]--

		simd := SliceContainsV0(s, target)
		simple := SliceContainsV1(s, target)
		cgoResult := cgo.SliceContains(s, target)

		require.False(t, simd)
		require.False(t, simple)
		require.False(t, cgoResult)
	})
}

// 16 alignment
func getData() ([]uint8, uint8) {
	s := make([]uint8, 1_000_000)

	for i := 0; i < len(s); i++ {
		s[i] = rand.N[uint8](5)
	}

	// always last match
	s[len(s)-1] = 10

	return s, 10
}
