#include "textflag.h"

// func vectorFloatAdditionV1(first, second, dst []float32)
TEXT ·vectorFloatAdditionV1(SB), NOSPLIT, $0
    LDP first_base+0(FP), (R0, R1)
    LDP second_base+24(FP), (R2, R3)
    LDP dst_base+48(FP), (R4, R5)

    MOVD $0, R7

loop:
    CMP R5, R7
    BGE done

    VLD1 (R0), [V0.S4]
    VLD1 (R2), [V1.S4]

    WORD $0x4e21d400 // fadd.4s v0, v0, v1

    VST1 [V0.S4], (R4)

    ADD $4, R7

    ADD $16, R4
    ADD $16, R0
    ADD $16, R2

    B loop

done:
    RET
