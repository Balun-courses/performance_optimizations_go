package main

import (
	"bufio"
	"context"
	"hw1/internal/server"
	"hw1/models"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

// TODO: add more log information for students
// panic just for stack trace

func main() {
	ctx, cancel := getSignalContext(context.Background())
	defer cancel()

	var (
		// we want to use bufio for concurrent processing, just extra task in the future
		input  = bufio.NewReaderSize(os.Stdin, models.MaxActionSize)
		output = bufio.NewWriter(os.Stdout)
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	)

	defer func() {
		err := output.Flush()

		if err != nil {
			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"can not flush output",
				slog.Any("error", err),
			)
		}
	}()

	ts := &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConnsPerHost: runtime.GOMAXPROCS(0) + 1,
	}

	srv := server.NewServer(
		input,
		output,
		logger,
		&http.Client{
			Transport: ts,
		},
	)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		// we know about pipe here
		err := os.Stdin.Close()

		if err != nil {
			panic(err)
		}
	}()

	srv.ListenAndServe(ctx)
	cancel()

	wg.Wait()
}
