package main

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

type MySliceHeader struct {
	data     *byte
	length   int
	capacity int
}

func (m *MySliceHeader) Len() int {
	return m.length
}

func (m *MySliceHeader) Cap() int {
	return m.capacity
}

var tmpBuffer = [8]byte{}

func (m *MySliceHeader) Push(a int64) {
	// check len and cap
	binary.LittleEndian.PutUint64(tmpBuffer[:], uint64(a))

	for i := 0; i < len(tmpBuffer); i++ {
		*m.getLast() = tmpBuffer[i]
		m.length++
	}
}

func (m *MySliceHeader) getLast() *byte {
	return (*byte)(unsafe.Add(unsafe.Pointer(m.data), m.length))
}

func (m *MySliceHeader) Pop() int64 {
	// check for length, shrink 1/4

	m.length -= 8
	for i := 0; i < 8; i++ {
		tmpBuffer[i] = *(*byte)(unsafe.Add(unsafe.Pointer(m.getLast()), i))
	}

	return int64(binary.LittleEndian.Uint64(tmpBuffer[:]))
}

func main() {
	s := &MySliceHeader{
		data:     myMalloc(100),
		length:   0,
		capacity: 100,
	}

	s.Push(1)
	s.Push(-2)
	s.Push(3)

	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
}

func myMalloc(size int) *byte {
	return unsafe.SliceData(make([]byte, size))
}
