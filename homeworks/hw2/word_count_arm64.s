#include "textflag.h"

#define ZERO(r) \
    MOVD $0, R2

// Здесь 32 - размер стек фрейма, 24 - размер аргументов
TEXT ·WordCount(SB), $32-24
    LDP data_base+0(FP), (R0, R1)
    MOVW $0, R6
    MOVW $0, R7

loop:
    CBZ R1, ret

    MOVW (R0), R5
    ADD $4, R0
    SUB $1, R1

    // Перед вызовом другой функции кладём данные и аргументы на стек
    MOVD R0, f-0(SP)
    MOVD R1, f-8(SP)
    MOVW R6, f-16(SP)
    MOVW R7, f-20(SP)

    MOVW R5, r-32(SP)

    // Обертка над unicode.IsSpace
    CALL ·isSpace(SB)
    MOVW s-24(SP), R8

    // Достаем результат и данные обратно
    MOVD f-0(SP), R0
    MOVD f-8(SP), R1
    MOVW f-16(SP), R6
    MOVW f-20(SP), R7

    // Test Bit and Branch if Zero
    TBZ	$0, R8, character
    MOVW $0, R6

    B loop

character:
    // R6 - word flag
    // Инвертируем все биты и делаем AND с 1
    MVN R6, R6
    AND $1, R6
    ADD R6, R7
    MOVW $1, R6
    B loop

ret:
    MOVW R7, ret+24(FP)
    RET
