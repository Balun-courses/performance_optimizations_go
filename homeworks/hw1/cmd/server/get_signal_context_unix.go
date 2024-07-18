//go:build unix

package main

import (
	"context"
	"os/signal"
	"syscall"
)

func getSignalContext(parentCtx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parentCtx, syscall.SIGALRM)
}
