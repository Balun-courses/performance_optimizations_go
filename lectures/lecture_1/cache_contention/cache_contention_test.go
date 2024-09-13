package cache_contention

import (
	"golang.org/x/sys/cpu"
	"math/rand/v2"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

const cacheLineSize = unsafe.Sizeof(cpu.CacheLinePad{})

func BenchmarkCacheContentionWithoutPadding(b *testing.B) {
	type WorkerStatistics struct {
		Value atomic.Int64
		_     [cacheLineSize - 8]byte
	}

	workers := runtime.NumCPU()
	statistics := make([]WorkerStatistics, workers)

	wg := new(sync.WaitGroup)

	for workerID := 0; workerID < workers; workerID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < b.N; j++ {
				statistics[workerID].Value.Add(externalStatisticSourceByID(int64(workerID)))
			}
		}()
	}

	wg.Wait()
}

func BenchmarkCacheContentionWithPadding(b *testing.B) {
	type WorkerStatistics struct {
		Value atomic.Int64
	}

	workers := runtime.NumCPU()
	statistics := make([]WorkerStatistics, workers)
	wg := new(sync.WaitGroup)

	for workerID := 0; workerID < workers; workerID++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < b.N; j++ {
				statistics[workerID].Value.Add(externalStatisticSourceByID(int64(workerID)))
			}
		}()
	}

	wg.Wait()
}

func externalStatisticSourceByID(id int64) int64 {
	return rand.N[int64](42*id%10 + 1)
}
