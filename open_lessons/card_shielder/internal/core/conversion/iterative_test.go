//go:build iterative

package conversion

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIterativeFibonacci(t *testing.T) {
	t.Parallel()

	require.Equal(t, uint64(12586269025), Fibonacci(50))
}

func BenchmarkFibonacciIterative(b *testing.B) {
	for i := 1; i <= 100_000; i++ {
		Fibonacci(i)
	}
}
