#include "textflag.h"

// Just for example
#define ZERO(r) \
    MOVD $0, r

// func SumSlice(s []int32) int64
TEXT ·SumSlice(SB), NOSPLIT, $0
    // Header слайса 24 байта
    // R0 - указатель на данные, R1 - длина
	LDP	slice_base+0(FP), (R0, R1)
    ZERO(R2)

loop:
    CBZ R1, done // Если R1 == 0, то переходим на метку done
    MOVW (R0), R9 // R9 = s[i]
    ADD R9, R2 // R2 += R9
    ADD $4, R0 // i++
    SUB $1, R1 // len--
    B loop // for range ...

done:
    MOVD R2, ret+24(FP)
    RET
