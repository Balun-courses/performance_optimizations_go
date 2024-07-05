package main

import (
	"github.com/stretchr/testify/require"
	"math/rand/v2"
	"strings"
	"testing"
)

func BenchmarkMatch(b *testing.B) {
	b.Run("fast match", func(b *testing.B) {
		b.StopTimer()
		data := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			fastMatch(data)
		}
	})

	b.Run("slow match", func(b *testing.B) {
		b.StopTimer()
		data := getData()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			slowMatch(data)
		}
	})
}

func getData() []string {
	data := make([]string, 0, 1<<21)

	for i := 0; i < cap(data); i++ {
		// see
		// https://github.com/golang/go/blob/master/src/runtime/string.go#L176
		data = append(data, strings.Repeat("1", 33<<6+1<<8+rand.N(2)))
	}

	return data
}

func TestMatch(t *testing.T) {
	data := getData()
	require.Equal(t, fastMatch(data), slowMatch(data))
}
