package main

import (
	"testing"
)

type Data struct {
	pointer *int
}

func BenchmarkEscapeAnalysis(b *testing.B) {
	b.Run("direction assignment", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			var number int
			data := &Data{
				pointer: &number,
			}
			_ = data
		}
	})

	b.Run("indirection assignment", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			var number int
			data := &Data{}
			data.pointer = &number
			_ = data
		}
	})
}
