package server

import (
	"bytes"
	"context"
	"hw1/models"
	"io"
	"log/slog"
	"net/http"
)

func (s *serverImpl) doRequestOp(ctx context.Context, action *models.DoRequestAction) {
	request, err := http.NewRequest(action.Method, action.Url.String(), bytes.NewBuffer(action.Body))

	if err != nil {
		s.invariantError(ctx, err)
		return
	}

	response, err := s.client.Do(request)

	defer func() {
		if response != nil && response.Body != nil {
			err = response.Body.Close()

			if err != nil {
				s.logger.ErrorContext(ctx, "can not close response body", slog.Any("error", err))
			}
		}
	}()

	if err != nil {
		s.sendResultMessage(ctx, models.DoRequestsOperation, models.DoRequestActionResult{
			RequestID: action.RequestID,
			Error:     err.Error(),
		})
		return
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		s.sendResultMessage(ctx, models.DoRequestsOperation, models.DoRequestActionResult{
			RequestID: action.RequestID,
			Error:     err.Error(),
		})
		return
	}

	s.sendResultMessage(ctx, models.DoRequestsOperation, models.DoRequestActionResult{
		RequestID: action.RequestID,
		Response:  body,
	})
}
