package conversion

import "runtime/debug"

func RecursiveFibonacci(n int) uint64 {
	if n <= 2 {
		return 1
	}

	return RecursiveFibonacci(n-1) + RecursiveFibonacci(n-2)
}

func OptimizedRecursiveFibonacci(n int) uint64 {
	debug.SetMaxStack(2000000000) // with optional demo

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

func IterativeFibonacci(n int) uint64 {
	result := make([]uint64, max(n, 2))
	result[0] = 1
	result[1] = 1

	for i := 2; i < n; i++ {
		result[i] = result[i-1] + result[i-2]
	}

	return result[n-1]
}

func OptimizedIterativeFibonacci(n int) uint64 {
	var (
		prev uint64
		cur  uint64
	)

	prev, cur = 1, 1

	for i := 2; i < n; i++ {
		prev, cur = cur, prev+cur
	}

	return cur
}
