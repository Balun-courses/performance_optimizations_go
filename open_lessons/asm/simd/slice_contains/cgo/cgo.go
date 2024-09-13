package cgo

/*
#include <stdio.h>
#include <stdint.h>
#include <arm_neon.h>
#include <stdbool.h>


bool slice_contains(const uint8_t *slice, size_t size, uint8_t value) {
    uint8x16_t val_vec = vdupq_n_u8(value);

    for (size_t i = 0; i < size; i += 16) {
        if (i + 16 > size) {
            break;
        }

        uint8x16_t data_vec = vld1q_u8(&slice[i]);
        uint8x16_t result_vec = vceqq_u8(data_vec, val_vec);
        uint16_t result = vaddvq_u8(result_vec);

        if (result) {
            return true;
        }
    }

    return false;
}
*/
import "C"

import (
	"unsafe"
)

func SliceContains(data []uint8, target uint8) bool {
	return bool(C.slice_contains((*C.uint8_t)(unsafe.SliceData(data)), C.size_t(len(data)), C.uint8_t(target)))
}
