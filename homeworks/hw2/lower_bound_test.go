package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"runtime"
	"testing"
)

func LowerBound(slice []int64, value int64) int64

func TestLowerBound(t *testing.T) {
	t.Parallel()

	fmt.Println("CURRENT ARCH")
	fmt.Println(runtime.GOARCH)
	type testCases struct {
		name   string
		arg    []int64
		value  int64
		result int64
	}

	tableTests := []testCases{
		{
			name:   "exact match",
			arg:    []int64{1, 2, 3, 4},
			value:  3,
			result: 2,
		},
		{
			name:   "empty",
			arg:    []int64{},
			value:  100,
			result: -1,
		},
		{
			name:   "none",
			arg:    []int64{10, 20, 30},
			value:  5,
			result: -1,
		},
		{
			name:   "last",
			arg:    []int64{5, 6, 7},
			value:  11,
			result: 2,
		},
		{
			name:   "one",
			arg:    []int64{-1},
			value:  -1,
			result: 0,
		},
		{
			name:   "one",
			arg:    []int64{-1},
			value:  1,
			result: 0,
		},
		{
			name:   "first",
			arg:    []int64{5, 10, 15},
			value:  7,
			result: 0,
		},
		{
			name:   "lower match",
			arg:    []int64{1, 2, 6, 8},
			value:  7,
			result: 2,
		},
	}

	for _, tt := range tableTests {
		tt := tt // go 1.21-

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.result, LowerBound(tt.arg, tt.value))
		})
	}
}
