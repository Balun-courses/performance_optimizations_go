//go:build unix || darwin

package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"hw1/models"
	"os"
	"sync"
	"syscall"
	"testing"
)

func TestServerSignalHandling(t *testing.T) {
	type testCase struct {
		signal syscall.Signal
	}

	testCases := []testCase{
		{
			signal: syscall.SIGALRM,
		},
	}

	pid := os.Getpid()
	wg := &sync.WaitGroup{}

	testRequest := models.GoMaxProcAction{
		Value: 10,
	}

	data, err := json.Marshal(testRequest)
	require.NoError(t, err)

	for _, tt := range testCases {
		r1, w1, err := os.Pipe()
		require.NoError(t, err)

		r2, w2, err := os.Pipe()
		require.NoError(t, err)

		os.Stdin = r1
		os.Stdout = w2

		wg.Add(1)
		go func() {
			defer wg.Done()
			main()
		}()

		_, err = w1.Write([]byte{models.GoMaxProcOperation})
		require.NoError(t, err)

		_, err = w1.Write(data)
		require.NoError(t, err)

		_, err = w1.Write([]byte{models.ActionsDelimiter})
		require.NoError(t, err)

		_, err = r2.Read(make([]byte, 1))
		require.NoError(t, err)

		err = syscall.Kill(pid, tt.signal)
		require.NoError(t, err)

		wg.Wait()
	}
}
