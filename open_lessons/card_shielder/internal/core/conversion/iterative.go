//go:build iterative

package conversion

func Fibonacci(n int) uint64 {
	result := make([]uint64, max(n, 2))
	result[0] = 1
	result[1] = 1

	for i := 2; i < n; i++ {
		result[i] = result[i-1] + result[i-2]
	}

	return result[n-1]
}
