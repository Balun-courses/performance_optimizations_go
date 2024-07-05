package main

import (
	"card_shielder/internal/api"
	"card_shielder/internal/config"
	"card_shielder/internal/repository"
	"context"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cardRepository := repository.NewCardRepository()
	mux := api.NewCardApi(logger, cardRepository).GetServeMux()

	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)
	registerProfilerEndpoints(mux)

	cfg := config.NewConstantConfig()
	apiAddress, err := cfg.GetApiAddress()

	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"can not start application",
			slog.Any("error", err),
		)
	}

	server := &http.Server{
		Addr:    apiAddress.String(),
		Handler: mux,
	}

	go func() {
		defer cancel()
		if err = server.ListenAndServe(); err != nil {
			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"server serve error",
				slog.Any("error", err),
			)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err = server.Shutdown(shutdownCtx)

	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"server shutdown error",
			slog.Any("error", err),
		)
		return
	}
}

func registerProfilerEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("goroutine"))
}
