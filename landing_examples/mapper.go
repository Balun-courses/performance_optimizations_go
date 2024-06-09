package main

import (
	"unsafe"
)

func fastMatch(exampleData []string) [][]byte {
	result := make([][]byte, 0, len(exampleData))

	for i := 0; i < len(exampleData); i++ {
		data := unsafe.Slice(unsafe.StringData(exampleData[i]), len(exampleData[i]))

		if externalMatchFunction(data) {
			result = append(result, data)
		}
	}

	return result
}

func slowMatch(exampleData []string) [][]byte {
	result := make([][]byte, 0)

	for i := 0; i < len(exampleData); i++ {
		data := []byte(exampleData[i])

		if externalMatchFunction(data) {
			result = append(result, data)
		}
	}

	return result
}

func externalMatchFunction(data []byte) bool {
	if len(data)%2 == 0 {
		return true
	}

	return false
}
