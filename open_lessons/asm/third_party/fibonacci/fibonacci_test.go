package fibonacci

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Fibonacci(n uint64) uint64

func TestFibonacci(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		n         uint64
		expResult uint64
	}{
		{
			name:      "zero",
			n:         0,
			expResult: 0,
		},
		{
			name:      "one",
			n:         1,
			expResult: 1,
		},
		{
			name:      "two",
			n:         2,
			expResult: 1,
		},
		{
			name:      "three",
			n:         3,
			expResult: 2,
		},
		{
			name:      "ten",
			n:         10,
			expResult: 55,
		},
		{
			name:      "70",
			n:         70,
			expResult: 190392490709135,
		},
	}

	for _, tt := range testCases {
		tt := tt // go 1.21-

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expResult, Fibonacci(tt.n))
		})
	}
}
