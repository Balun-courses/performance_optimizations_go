package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const count = 100_000_000

func main() {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(wd, "/lections/lection_3/data.txt"))

	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(f)

	for i := 0; i <= count; i++ {
		var chunk string

		if i == count {
			chunk = generateChunk()
		} else {
			chunk = generateChunk() + " "
		}

		_, err := writer.WriteString(chunk)

		if err != nil {
			panic(err)
		}
	}
}

func generateChunk() string {
	return strings.Repeat("data", 42)
}
