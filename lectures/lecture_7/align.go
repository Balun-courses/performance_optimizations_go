package main

import (
	"fmt"
	"unsafe"
)

type T1 struct {
	a int8
	b int64
	c int16
}

type T2 struct {
	a int8
	c int16
	b int64
}

func main() {
	fmt.Println(unsafe.Sizeof(T1{})) // 24
	fmt.Println(unsafe.Sizeof(T2{})) // 16

	fmt.Println(unsafe.Alignof(T1{}.a)) // 1
	fmt.Println(unsafe.Alignof(T1{}.b)) // 8
	fmt.Println(unsafe.Alignof(T1{}.c)) // 2
}
