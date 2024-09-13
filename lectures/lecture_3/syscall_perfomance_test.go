package main

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkSyscallReadPerformance(b *testing.B) {
	wd, err := os.Getwd()
	require.NoError(b, err)

	dataPath := filepath.Join(wd, "data.txt")
	constChunk := generateChunk()

	b.Run("buffered", func(b *testing.B) {
		f, err := os.Open(dataPath)
		require.NoError(b, err)

		reader := bufio.NewReader(f)

		for i := 0; i < b.N; i++ {
			data, err := reader.ReadString(' ')
			require.NoError(b, err)
			require.Equal(b, constChunk, data[:len(data)-1])
		}
	})

	b.Run("not buffered", func(b *testing.B) {
		f, err := os.Open(dataPath)
		require.NoError(b, err)

		for i := 0; i < b.N; i++ {
			var data string

			_, err = fmt.Fscan(f, &data)
			require.NoError(b, err)
			require.Equal(b, constChunk, data)
		}
	})
}
