#include "textflag.h"

#define ZERO(r) \
    MOVD $0, R2

// func LowerBound(s []int64, v int64) int64
TEXT ·LowerBound(SB), NOSPLIT, $0
    LDP slice_base+0(FP), (R0, R1)
    MOVD value+24(FP), R2
    MOVD $-1, R3

loop:
    SUB R3, R1, R4
    CMP $1, R4
    BLE done

    MOVD $1, R8
    LSR R8, R4, R7
    ADD R3, R7
    MOVD R7, R13

    MOVD $8, R11
    MUL R11, R7, R7 // middle = (right - left)/2 + left
    ADD R7, R0, R12

    MOVD (R12), R5
    CMP R2, R5

    CSEL LE, R13, R3, R3
    CSEL GT, R13, R1, R1

    B loop

done:
    MOVD R3, ret+32(FP)
    RET
