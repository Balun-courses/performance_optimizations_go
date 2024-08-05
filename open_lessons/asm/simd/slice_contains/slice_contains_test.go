package slice_contains

import (
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
	b.Run("SIMD contains", func(b *testing.B) {
		b.StopTimer()
		s, target := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			SliceContainsV1(s, target)
		}
	})

	b.Run("simple contains", func(b *testing.B) {
		b.StopTimer()
		s, target := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			SliceContainsV0(s, target)
		}
	})
}

func TestSliceContains(t *testing.T) {
	s, target := getData()
	require.Equal(t, target, s[len(s)-1])

	// always last match
	require.False(t, slices.Contains(s[:len(s)-1], target))

	simd := SliceContainsV0(s, target)
	simple := SliceContainsV1(s, target)

	require.True(t, simd)
	require.True(t, simple)
}

// 16-bit alignment
func getData() ([]uint8, uint8) {
	s := make([]uint8, 1_000_000)

	for i := 0; i < len(s); i++ {
		s[i] = rand.N[uint8](5)
	}

	// always last match
	s[len(s)-1] = 10

	return s, 10
}
