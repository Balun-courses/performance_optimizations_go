//go:build recursive

package conversion

func Fibonacci(n int) uint64 {
	if n <= 2 {
		return 1
	}

	return Fibonacci(n-1) + Fibonacci(n-2)
}
