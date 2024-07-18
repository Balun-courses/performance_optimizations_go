//go:build windows

package main

import "context"

func getSignalContext(parentCtx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parentCtx)
}
