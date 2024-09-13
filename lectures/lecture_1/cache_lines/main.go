package main

import (
	"fmt"
	"golang.org/x/sys/cpu"
	"unsafe"
)

const cacheLineSize = unsafe.Sizeof(cpu.CacheLinePad{})

func main() {
	a := 42
	fmt.Println(fmt.Sprintf("%064b", &a))
	fmt.Println(cacheLineSize)
}
