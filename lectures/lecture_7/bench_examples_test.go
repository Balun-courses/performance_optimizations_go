package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

type data struct {
	a [10224323]int64
}

// go test -bench=. -benchmem
func BenchmarkSimple(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		creationFunc[data]()
	}
}

func creationFunc[T any]() *T {
	return new(T)
}

// go test -bench=. -benchmem -benchtime=100x
func BenchmarkWithTimer(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		var a int
		s := make([]int64, 1000)
		b.StartTimer()
		// or ResetTimer

		a += len(s)
	}
}

// go test -bench=. -benchmem -cpu=x
func BenchmarkParallel(b *testing.B) {
	// -cpu=x
	b.SetParallelism(runtime.NumCPU())
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		// no timer settings here

		for pb.Next() {
			// action
		}
	})
}

func BenchmarkWrongRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		recursiveFunc(b.N)
	}
}

func recursiveFunc(a int) {

}

func TestWrongComparison(t *testing.T) {
	res := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			recursiveFunc(i)
		}
	})
	if res.NsPerOp() <= 1000 {
	}
}

func BenchmarkWrongSingleAction(b *testing.B) {
	recursiveFunc(1_000_000)
}

func BenchmarkZeroCopyConverter(b *testing.B) {
	b.StopTimer()
	s := make([]byte, 10_000_000)
	b.StartTimer()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		zeroCopyConverter(s)
	}
}

func BenchmarkStrconv(b *testing.B) {
	b.Run("sprinf", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			fmt.Sprintf("%d", 42)
		}
	})

	b.Run("strconv", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			strconv.FormatInt(42, 10)
		}
	})
}

func BenchmarkStringBuilder(b *testing.B) {
	b.Run("stringBuilder", func(b *testing.B) {
		b.ReportAllocs()

		sb := &strings.Builder{}

		for i := 0; i < b.N; i++ {
			sb.WriteByte(byte(i))
		}

		result := sb.String()
		_ = result
	})

	b.Run("simple addition", func(b *testing.B) {
		b.ReportAllocs()

		var result string

		for i := 0; i < b.N; i++ {
			result += string(byte(i))
		}
	})
}

type TestData struct {
	a string
	b int64
	c int64
}

func BenchmarkLoopAllocation(b *testing.B) {

	b.Run("loop allocation", func(b *testing.B) {
		b.ReportAllocs()

		b.StopTimer()
		sd, err := json.Marshal(TestData{
			a: "213132",
			b: 123,
			c: 123124,
		})
		require.NoError(b, err)
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			d := Data{}

			err = json.Unmarshal(sd, &d)
			require.NoError(b, err)
		}
	})

	b.Run("outer loop allocation", func(b *testing.B) {
		b.ReportAllocs()

		b.StopTimer()
		sd, err := json.Marshal(TestData{
			a: "213132",
			b: 123,
			c: 123124,
		})
		require.NoError(b, err)
		b.StartTimer()

		d := Data{}
		for i := 0; i < b.N; i++ {
			err = json.Unmarshal(sd, &d)
			require.NoError(b, err)
		}
	})
}

func zeroCopyConverter(bigSlice []byte) []byte {
	str := unsafe.String(unsafe.SliceData(bigSlice), len(bigSlice))
	s := unsafe.Slice(unsafe.StringData(str), len(str))

	return s
}
