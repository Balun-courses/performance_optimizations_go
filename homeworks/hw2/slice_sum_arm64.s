#include "textflag.h"

#define ZERO(r) \
    MOVD $0, R2

// TODO: comments

// func SumSlice(s []int32) int64
TEXT ·SumSlice(SB), NOSPLIT, $0
	LDP	x_base+0(FP), (R0, R1)	// R0 = s.ptr, R1 = len(s)
    ZERO(R2)

loop:
    CBZ R1, done
    MOVW (R0), R9
    ADD R9, R2
    ADD $4, R0
    SUB $1, R1
    B loop

done:
    MOVD R2, ret+24(FP)
    RET
