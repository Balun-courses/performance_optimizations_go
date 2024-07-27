#include "textflag.h"

// I: n >= 0
// func Fibonacci(n uint64) uint64 {
//	if n == 0 {
//		return 0
//	}
//
//	var (
//		prev uint64 = 0
//		cur  uint64 = 1
//	)
//
//	for i := 2; i <= int(n); i++ {
//		prev, cur = cur, prev+cur
//	}
//
//	return cur
//}

// func Fibonacci(n uint64) uint64
TEXT ·Fibonacci(SB), NOSPLIT, $0
    MOVD n+0(FP), R0 // R0 := n

    MOVD $0, R2 // R2 := 0
    CBZ R0, done // if R0 == 0 { return R2}

    MOVD $0, R1 // R1 := 0
    MOVD $1, R2 // R2 := 1

    MOVD $2, R3 // R3 := 2
loop:
    CMP R0, R3 // for R3 <= R0 {}
    BGT done

    MOVD R2, R7 // R7 := R2
    ADD R1, R2 // R2 += R1
    MOVD R7, R1 // R1 = R7

    ADD $1, R3 // R3++

    B loop
done:
    MOVD R2, ret+8(FP)
    RET
