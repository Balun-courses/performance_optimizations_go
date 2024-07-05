#include "textflag.h"

#define ZERO(r) \
    MOVD $0, R2

// TODO: comments


// func WordCount(data []rune) int32
// R5 - cur symbol
// R6 - word flag
// R7 - result
TEXT ·WordCount(SB), $32-24
    LDP data_base+0(FP), (R0, R1) // R0 - data_ptr, R1 - len
    MOVW $0, R6
    MOVW $0, R7

loop:
    CBZ R1, ret

    MOVW (R0), R5
    ADD $4, R0
    SUB $1, R1

    MOVD R0, f-0(SP)
    MOVD R1, f-8(SP)
    MOVW R6, f-16(SP)
    MOVW R7, f-20(SP)

    MOVW R5, r-32(SP)
    CALL ·isSpace(SB)
    MOVW s-24(SP), R8

    MOVD f-0(SP), R0
    MOVD f-8(SP), R1
    MOVW f-16(SP), R6
    MOVW f-20(SP), R7

    TBZ	$0, R8, character
    MOVW $0, R6

    B loop

character:
    MVN R6, R6
    AND $1, R6
    ADD R6, R7
    MOVW $1, R6
    B loop

ret:
    MOVW R7, ret+24(FP)
    RET
