#include "textflag.h"

#define ZERO(r) \
    MOVD $0, R2

// func LowerBound(data []int64, target int64) int64
TEXT ·LowerBound(SB), NOSPLIT, $0
    LDP slice_base+0(FP), (R0, R1) // R0 := data.ptr, R1 := len(data)
    MOVD value+24(FP), R2 // R2 := target
    MOVD $-1, R3 // R3 := -1

loop:
    // for right - left > 1
    SUB R3, R1, R4 // R4 = R1 - R3
    CMP $1, R4 // comparison
    BLE done

    // middle = (right - left)/2 + left
    MOVD $1, R8 // R8 := 1
    LSR R8, R4, R7 // R7 := R4 >> R8
    ADD R3, R7 // R7 += R3
    MOVD R7, R13

    MOVD $8, R11
    MUL R11, R7, R7
    ADD R7, R0, R12
    // R5 := *(data.ptr + 8 * middle)
    MOVD (R12), R5

    CMP R2, R5

    // if R5 <= target {
    //     R3 = middle
    // } else {
    //     R1 = middle
    //}
    CSEL LE, R13, R3, R3
    CSEL GT, R13, R1, R1

    B loop

done:
    MOVD R3, ret+32(FP)
    RET
