//go:build optimized_iterative

package conversion

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOptimizedIterativeFibonacci(t *testing.T) {
	t.Parallel()

	require.Equal(t, uint64(12586269025), Fibonacci(50))
}

func BenchmarkFibonacciOptimizedIterative(b *testing.B) {
	for i := 1; i <= 10_000; i++ {
		Fibonacci(i)
	}
}
