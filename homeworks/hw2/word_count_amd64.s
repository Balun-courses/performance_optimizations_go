#include "textflag.h"

#define ZERO(r) \
    MOVQ $0, R2

// func WordCount(data []rune) int32
TEXT ·WordCount(SB), $32-24
    MOVQ data_base+0(FP), AX
    MOVQ len+8(FP), R9
    MOVL $0, DX
    MOVL $0, CX

loop:
    CMPQ R9, $0
    JE ret

    MOVLQSX (AX), R11
    ADDQ $4, AX
    SUBQ $1, R9

    MOVQ AX, f-0(SP)
    MOVQ R9, f-8(SP)
    MOVL DX, f-16(SP)
    MOVL CX, f-20(SP)

    MOVL R11, r-32(SP)
    CALL ·isSpace(SB)
    MOVQ s-24(SP), R8
    MOVQ f-0(SP), AX
    MOVQ f-8(SP), R9
    MOVL f-16(SP), DX // no negative here
    MOVL f-20(SP), CX

    BTQ $0, R8
    JNC character
    MOVL $0, DX

    JMP loop

character:
    XORQ $0xffffffff, DX

    ANDQ $1, DX
    ADDQ DX, CX
    MOVL $1, DX
    JMP loop

ret:
    MOVL CX, ret+24(FP)
    RET
