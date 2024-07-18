package main

import (
	"bufio"
	"bytes"
	"cmp"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"hw1/models"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"testing"
	"time"
)

// TODO: timeout
func TestGoMacProcOp(t *testing.T) {
	type testCase struct {
		name      string
		action    models.GoMaxProcAction
		state     func()
		expResult models.GoMaxProcActionResult
	}

	testCases := []testCase{
		{
			name: "change action",
			action: models.GoMaxProcAction{
				Value: 10,
			},
			state: func() {
				runtime.GOMAXPROCS(32)
			},
			expResult: models.GoMaxProcActionResult{
				PreviousValue: 32,
			},
		},
		{
			name: "no change action",
			action: models.GoMaxProcAction{
				Value: 0,
			},
			state: func() {
				runtime.GOMAXPROCS(64)
			},
			expResult: models.GoMaxProcActionResult{
				PreviousValue: 64,
			},
		},
	}

	for _, tt := range testCases {
		tt.state()

		r1, w1, err := os.Pipe()
		require.NoError(t, err)

		r2, w2, err := os.Pipe()
		require.NoError(t, err)

		// check windows support, extra files is not supported on it
		os.Stdin = r1
		os.Stdout = w2

		serialized, err := json.Marshal(tt.action)
		require.NoError(t, err)

		_, err = w1.Write([]byte{models.GoMaxProcOperation})
		require.NoError(t, err)

		_, err = w1.Write(serialized)
		require.NoError(t, err)

		_, err = w1.Write([]byte{models.ActionsDelimiter})
		require.NoError(t, err)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			main()
		}()

		err = r2.SetReadDeadline(time.Now().Add(time.Second * 5))
		data, err := bufio.NewReader(r2).ReadSlice(models.ActionsDelimiter)
		require.NoError(t, err)
		require.Equal(t, models.GoMaxProcOperation, data[0])

		var result models.GoMaxProcActionResult
		err = json.Unmarshal(data[1:len(data)-1], &result)
		require.NoError(t, err)

		require.Equal(t, tt.expResult, result)
		require.Equal(t, cmp.Or(tt.action.Value, tt.expResult.PreviousValue), runtime.GOMAXPROCS(-1))

		err = w1.Close()
		require.NoError(t, err)

		wg.Wait()
	}
}

func TestServerDoRequestOp(t *testing.T) {
	body := bytes.Repeat([]byte{'x'}, models.MaxActionSize-0xfff)

	var handler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		require.Equal(t, request.Method, http.MethodPost)

		data, err := io.ReadAll(request.Body)
		require.NoError(t, err)

		defer func() {
			err = request.Body.Close()
			require.NoError(t, err)
		}()

		require.Equal(t, body, data)

		_, err = writer.Write([]byte("123"))
		require.NoError(t, err)
	}

	testServer := httptest.NewServer(handler)
	t.Cleanup(func() {
		testServer.Close()
	})

	serverUrl, err := url.Parse(testServer.URL)
	require.NoError(t, err)

	r1, w1, err := os.Pipe()
	require.NoError(t, err)

	r2, w2, err := os.Pipe()
	require.NoError(t, err)

	os.Stdin = r1
	os.Stdout = w2

	serialized, err := json.Marshal(models.DoRequestAction{
		RequestID: 123,
		Url:       serverUrl,
		Method:    http.MethodPost,
		Body:      body,
	})

	require.NoError(t, err)

	_, err = w1.Write([]byte{models.DoRequestsOperation})
	require.NoError(t, err)

	_, err = w1.Write(serialized)
	require.NoError(t, err)

	_, err = w1.Write([]byte{models.ActionsDelimiter})
	require.NoError(t, err)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		main()
	}()

	err = r2.SetReadDeadline(time.Now().Add(time.Second * 5))
	require.NoError(t, err)

	data, err := bufio.NewReader(r2).ReadSlice(models.ActionsDelimiter)
	require.NoError(t, err)
	require.Equal(t, data[0], models.DoRequestsOperation)

	var res models.DoRequestActionResult
	err = json.Unmarshal(data[1:len(data)-1], &res)
	require.NoError(t, err)

	require.Equal(t, models.DoRequestActionResult{
		RequestID: 123,
		Response:  []byte("123"),
	}, res)

	err = w1.Close()
	require.NoError(t, err)

	wg.Wait()
}

func TestSetMemoryLimitOp(t *testing.T) {
	type testCase struct {
		name      string
		action    models.SetMemoryLimitAction
		state     func()
		expResult models.SetMemoryLimitActionResult
	}

	testCases := []testCase{
		{
			name: "change action",
			action: models.SetMemoryLimitAction{
				Value: 256_000,
			},
			state: func() {
				debug.SetMemoryLimit(123_000)
			},
			expResult: models.SetMemoryLimitActionResult{
				PreviousValue: 123_000,
			},
		},
		{
			name: "no change action",
			action: models.SetMemoryLimitAction{
				Value: -1,
			},
			state: func() {
				debug.SetMemoryLimit(123_000)
			},
			expResult: models.SetMemoryLimitActionResult{
				PreviousValue: 123_000,
			},
		},
	}

	for _, tt := range testCases {
		tt.state()

		r1, w1, err := os.Pipe()
		require.NoError(t, err)

		r2, w2, err := os.Pipe()
		require.NoError(t, err)

		os.Stdin = r1
		os.Stdout = w2

		serialized, err := json.Marshal(tt.action)
		require.NoError(t, err)

		_, err = w1.Write([]byte{models.SetMemoryLimitOperation})
		require.NoError(t, err)

		_, err = w1.Write(serialized)
		require.NoError(t, err)

		_, err = w1.Write([]byte{models.ActionsDelimiter})
		require.NoError(t, err)

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			main()
		}()

		err = r2.SetReadDeadline(time.Now().Add(time.Second * 5))
		require.NoError(t, err)

		data, err := bufio.NewReader(r2).ReadSlice(models.ActionsDelimiter)
		require.NoError(t, err)
		require.Equal(t, models.SetMemoryLimitOperation, data[0])

		var res models.SetMemoryLimitActionResult

		err = json.Unmarshal(data[1:len(data)-1], &res)
		require.NoError(t, err)
		require.Equal(t, tt.expResult, res)

		if tt.action.Value >= 0 {
			require.Equal(t, tt.action.Value, debug.SetMemoryLimit(-1))
		}

		err = w1.Close()
		require.NoError(t, err)

		wg.Wait()
	}
}
