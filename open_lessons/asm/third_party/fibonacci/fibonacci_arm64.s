#include "textflag.h"

// TODO: add code emulation

// func Fibonacci(n uint64) uint64
TEXT ·Fibonacci(SB), NOSPLIT, $0
    MOVD n+0(FP), R0

    MOVD $0, R2
    CBZ R0, done

    MOVD $0, R1
    MOVD $1, R2

    MOVD $2, R3
loop:
    CMP R0, R3
    BGT done

    MOVD R2, R7
    ADD R1, R2
    MOVD R7, R1

    ADD $1, R3

    B loop
done:
    MOVD R2, ret+8(FP)
    RET
