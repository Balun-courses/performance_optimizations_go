#include "textflag.h"

// func vectorAdditionV1(first, second, dst []uint8)
TEXT ·vectorAdditionV1(SB), NOSPLIT, $0
    LDP first_base+0(FP), (R0, R1)
    LDP second_base+24(FP), (R2, R3)
    LDP dst_base+48(FP), (R4, R5)

    MOVD $0, R7
    MOVD R4, R11
    MOVD R0, R12
    MOVD R2, R13

loop:
    CMP R5, R7
    BGE done

    VLD1 (R12), [V1.B16]
    VLD1 (R13), [V2.B16]
    VADD V1.B16, V2.B16, V3.B16

    VST1 [V3.B16], (R11)

    ADD $16, R7

    ADD $16, R11
    ADD $16, R12
    ADD $16, R13

    B loop

done:
    RET
