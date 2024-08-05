#include "textflag.h"

#define ZERO(r) \
    MOVD $0, R2

// TODO: comments

// func LowerBound(s []int64, v int64) int64
// R3 - leftBorder, R1 - rightBorder, R2 - target, R4 - middle
TEXT ·LowerBound(SB), NOSPLIT, $0
    LDP slice_base+0(FP), (R0, R1) // R0 = s.ptr, R1 = len(s)
    MOVD value+24(FP), R2 // R2 = v
    MOVD $-1, R3 // R3 = -1

loop:
    SUB R3, R1, R4 // right - left
    CMP $1, R4 // right - left > 1
    BLE done

    MOVD $1, R8
    LSR R8, R4, R7 // / 2
    ADD R3, R7
    MOVD R7, R13

    MOVD $8, R11
    MUL R11, R7, R7
    ADD R7, R0, R12

    MOVD (R12), R5
    CMP R2, R5

// assembly
  // CMP R1, R3        // Сравнивает значения в регистрах R1 и R3.
  // CSEL R6, R1, R3, LT  // R6 = если R1 < R3, то R1, иначе R3
    CSEL LE, R13, R3, R3
    CSEL GT, R13, R1, R1

    B loop

done:
    MOVD R3, ret+32(FP)
    RET
