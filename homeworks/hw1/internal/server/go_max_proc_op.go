package server

import (
	"context"
	"hw1/models"
	"runtime"
)

func (s *serverImpl) goMaxProcOp(ctx context.Context, action *models.GoMaxProcAction) {
	s.sendResultMessage(ctx, models.GoMaxProcOperation, models.GoMaxProcActionResult{
		PreviousValue: runtime.GOMAXPROCS(action.Value),
	})
}
