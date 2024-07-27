//go:build optimized_recursive

package conversion

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOptimizedRecursiveFibonacci(t *testing.T) {
	t.Parallel()

	require.Equal(t, uint64(12586269025), Fibonacci(50))
}

func BenchmarkFibonacciOptimizedRecursive(b *testing.B) {
	for i := 1; i <= 10_000; i++ {
		Fibonacci(i)
	}
}
