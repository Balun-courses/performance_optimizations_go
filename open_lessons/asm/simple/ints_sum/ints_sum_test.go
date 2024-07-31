package ints_sum

import (
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func IntSum(a, b int64) int64

func TestIntSum(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		first     int64
		second    int64
		expResult int64
	}{
		{
			name:      "small positive",
			first:     1,
			second:    2,
			expResult: 3,
		},
		{
			name:      "small negative",
			first:     -1,
			second:    -2,
			expResult: -3,
		},
		{
			name:      "mixed negative",
			first:     -3,
			second:    2,
			expResult: -1,
		},
		{
			name:      "big positive",
			first:     math.MaxInt32,
			second:    math.MaxInt32,
			expResult: math.MaxInt32 * 2,
		},
		{
			name:      "big negative",
			first:     math.MinInt32,
			second:    math.MinInt32,
			expResult: math.MinInt32 * 2,
		},
		{
			name:      "check threshold",
			first:     math.MaxInt64,
			second:    -1,
			expResult: math.MaxInt64 - 1,
		},
	}

	for _, tt := range testCases {
		tt := tt // go 1.21-

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.expResult, IntSum(tt.first, tt.second))
		})
	}
}
