//go:build optimized_recursive

package conversion

func Fibonacci(n int) uint64 {
	mem := make([]uint64, n+1)

	var fib func(n int) uint64

	fib = func(n int) uint64 {
		if mem[n] != 0 {
			return mem[n]
		}

		if n <= 2 {
			mem[n] = 1
		} else {
			mem[n] = fib(n-1) + fib(n-2)
		}

		return mem[n]
	}

	return fib(n)
}
