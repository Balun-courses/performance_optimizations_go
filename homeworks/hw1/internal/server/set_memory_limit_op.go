package server

import (
	"context"
	"hw1/models"
	"runtime/debug"
)

func (s *serverImpl) setMemoryLimitOp(ctx context.Context, action *models.SetMemoryLimitAction) {
	s.sendResultMessage(ctx, models.SetMemoryLimitOperation, models.SetMemoryLimitActionResult{
		PreviousValue: debug.SetMemoryLimit(action.Value),
	})
}
