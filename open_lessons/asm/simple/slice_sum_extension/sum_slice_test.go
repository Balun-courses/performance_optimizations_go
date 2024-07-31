package main

import (
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func SumSlice(x []int32) int64

func TestSumSlice(t *testing.T) {
	t.Parallel()

	type testCases struct {
		name   string
		arg    []int32
		result int64
	}

	tableTests := []testCases{
		{
			name:   "positive",
			arg:    []int32{1, 2, 3},
			result: 6,
		},
		{
			name:   "negative",
			arg:    []int32{-1, -2, -3},
			result: -6,
		},
		{
			name:   "mixed",
			arg:    []int32{-1, 2, -3},
			result: -2,
		},
		{
			name:   "int32 threshold",
			arg:    []int32{math.MaxInt32, 1},
			result: 2147483648,
		},
		{
			name:   "int32 threshold",
			arg:    []int32{math.MinInt32, -1},
			result: -2147483649,
		},
	}

	for _, tt := range tableTests {
		tt := tt // go 1.21-
		
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.result, SumSlice(tt.arg))
		})
	}
}
