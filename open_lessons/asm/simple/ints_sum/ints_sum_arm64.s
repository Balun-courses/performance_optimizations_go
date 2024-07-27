#include "textflag.h"

// "textflag.h" содержит стандартные директивы, например, NOSPLIT

// func IntSum(a, b int64) int64
// . и / заменены из-за синтаксиса, например, (math∕rand·Int)

TEXT ·IntSum(SB), NOSPLIT, $0 // $x-y (x - размер фрейма, y - размер аргументов. В случае NOSPLIT можно $0)
    LDP slice_base+0(FP), (R0, R1) // Загружаем в R0 и R1 аргументы функции
    ADD R0, R1 // Прибавляем к R1 += R0
    MOVD R1, ret+16(FP) // R1 to return
    RET // return
