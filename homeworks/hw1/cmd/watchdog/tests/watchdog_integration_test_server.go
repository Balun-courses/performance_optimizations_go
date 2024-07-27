package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"hw1/models"
	"log/slog"
	"net/http"
	"os"
)

// TODO: add more log information for students
// panic just for stack trace

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	logger.Info("command line arguments", slog.Any("args", os.Args))

	addr := flag.String("test_server_result_url", "", "test_server_result_url")
	flag.Parse()

	if *addr == "" {
		panic("got empty test_server_result_url") // just for stacktrace
	}

	testEnv := os.Getenv("TEST_ENV")

	if testEnv == "" {
		panic("expected not empty TEST_ENV")
	}

	if testEnv != "TEST_ENV_DATA" {
		panic("vector_addition inh mismatch")
	}

	reader := bufio.NewReader(os.Stdin)

	data, err := reader.ReadSlice(models.ActionsDelimiter)

	if err != nil {
		panic(err)
	}

	if data[0] != models.GoMaxProcOperation {
		panic("expected go mac proc command")
	}

	var request models.GoMaxProcAction
	err = json.Unmarshal(data[1:len(data)-1], &request)

	if request.Value != 32 {
		panic("expected value 32")
	}

	req, err := http.NewRequest(http.MethodPost, *addr, bytes.NewBuffer([]byte("OK")))

	if err != nil {
		panic(err)
	}

	_, err = http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	response, err := json.Marshal(models.GoMaxProcActionResult{
		PreviousValue: 13,
	})

	if err != nil {
		panic(err)
	}

	fmt.Fprint(os.Stdout, string(models.GoMaxProcOperation))
	fmt.Fprint(os.Stdout, string(response))
	fmt.Fprint(os.Stdout, string(models.ActionsDelimiter))

	fmt.Fprintln(os.Stderr, "SO, TRY TO LIVE")

	for i := 0; i < 1_000_000_000; i++ {
		liveRequest, err := http.NewRequest(http.MethodPost, *addr, bytes.NewBuffer([]byte("COUNTER")))

		if err != nil {
			panic(err)
		}

		_, err = http.DefaultClient.Do(liveRequest)

		if err != nil {
			panic(err)
		}
	}
}
