package conversion

import (
	"testing"
)

func BenchmarkFibonacciIterative(b *testing.B) {
	b.Run("Iterative", func(b *testing.B) {
		for i := 1; i <= 100000; i++ {
			IterativeFibonacci(i)
		}
	})

	b.Run("Optimized iterative", func(b *testing.B) {
		for i := 1; i <= 100000; i++ {
			OptimizedIterativeFibonacci(i)
		}
	})

	b.Run("Optimized recursion", func(b *testing.B) {
		for i := 1; i <= 35; i++ {
			RecursiveFibonacci(i)
		}
	})

	b.Run("Recursion", func(b *testing.B) {
		for i := 1; i < 35; i++ {
			RecursiveFibonacci(i)
		}
	})
}
