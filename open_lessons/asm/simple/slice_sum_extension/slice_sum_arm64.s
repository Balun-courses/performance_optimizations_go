#include "textflag.h"

// Just for example
#define ZERO(r) \
    MOVD $0, r

// type SliceHeader struct {
//  Data uintptr
//	Len  int
//	Cap  int
//}

// func SumSlice(x []int32) int64 {
//	var result int64
//
//	for i := range x {
//		result += int64(x[i])
//	}
//
//	return result
//}

// func SumSlice(s []int32) int64
TEXT ·SumSlice(SB), NOSPLIT, $0
    // Header слайса 24 байта
    // R0 - указатель на данные, R1 - длина
	LDP	slice_base+0(FP), (R0, R1)
    ZERO(R2)

loop:
    CBZ R1, done // for R1 != 0 {}
    MOVW (R0), R9 // R9 = s[i]
    ADD R9, R2 // R2 += R9
    ADD $4, R0 // R0++
    SUB $1, R1 // R1--
    B loop //

done:
    MOVD R2, ret+24(FP)
    RET
