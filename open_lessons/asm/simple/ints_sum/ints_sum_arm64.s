#include "textflag.h"

// func IntSum(a, b int64) int64
TEXT ·IntSum(SB), NOSPLIT, $0
    LDP slice_base+0(FP), (R0, R1)
    ADD R0, R1
    MOVD R1, ret+16(FP)
    RET
