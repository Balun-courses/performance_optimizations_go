#include "textflag.h"

//func SliceContainsV1(s []uint8, target uint8) bool
TEXT ·SliceContainsV1(SB), NOSPLIT, $0
    LDP slice_base+0(FP), (R0, R1)
    MOVB target+24(FP), R2
    VDUP R2, V1.B16

loop:
    CBZ R1, no

    VLD1.P 16(R0), [V2.B16]
    VCMEQ V1.B16, V2.B16, V3.B16

    VMOV V3.D[0], R4
    VMOV V3.D[1], R5

    CBNZ R4, yes
    CBNZ R5, yes

    SUB $16, R1

    B loop

no:
    MOVD $0, R5
    MOVD R5, ret+32(FP)
    RET


yes:
    MOVD $1, R5
    MOVD R5, ret+32(FP)
    RET
