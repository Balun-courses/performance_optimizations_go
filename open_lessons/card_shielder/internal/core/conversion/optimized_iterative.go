//go:build optimized_iterative

package conversion

func Fibonacci(n int) uint64 {
	var (
		prev uint64
		cur  uint64
	)

	prev, cur = 0, 1

	for i := 2; i <= n; i++ {
		prev, cur = cur, prev+cur
	}

	return cur
}
