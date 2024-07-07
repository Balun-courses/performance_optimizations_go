#include "textflag.h"

// TODO: comments

// func LowerBound(s []int64, v int64) int64
// R8 - leftBorder, R9 - rightBorder, CX - target, R10 - middle
TEXT ·LowerBound(SB), NOSPLIT, $0
    MOVQ slice_base+0(FP), AX // AX=ptr
    MOVQ len+8(FP), R9 // R9=len
    MOVD value+24(FP), CX // CX = v
    MOVD $-1, R8 // R8 = -1

loop:
    MOVQ R9, R10
    SUBQ R8, R10

    CMPQ R10, $1 // right - left > 1
    JLE done

    SHRQ $1, R10 // middle = (right - left)/2 + left
    ADDQ R8, R10 // overflow guard

    MOVQ 0(AX)(R10*8), R11

    CMPQ R11, CX

    CMOVQLE R10, R8
    CMOVQGT R10, R9

    JMP loop

done:
    MOVQ R8, ret+32(FP)
    RET
