package server

import (
	"context"
	"encoding/json"
	"fmt"
	"hw1/models"
	"io"
	"log/slog"
	"net/http"
)

// TODO: add more log information for students

type (
	reader interface {
		io.Reader
		ReadSlice(delim byte) (line []byte, err error)
	}

	writer interface {
		io.Writer
		WriteByte(c byte) error
		Flush() error
	}
)

type Server interface {
	ListenAndServe(ctx context.Context)
}

var _ Server = (*serverImpl)(nil)

type serverImpl struct {
	input  reader
	output writer
	logger *slog.Logger
	client *http.Client
}

func NewServer(
	input reader,
	output writer,
	logger *slog.Logger,
	client *http.Client,
) *serverImpl {
	return &serverImpl{
		input:  input,
		output: output,
		logger: logger,
		client: client,
	}
}

func (s *serverImpl) ListenAndServe(ctx context.Context) {
	defer func() {
		s.logger.Info("flush")
		s.output.Flush()
	}()

	for {
		if ctx.Err() != nil {
			return
		}

		data, err := s.input.ReadSlice(models.ActionsDelimiter)

		if err != nil && err != io.EOF {
			s.invariantError(ctx, err)
			return
		}

		if err == io.EOF && len(data) == 0 {
			return
		}

		err = s.processAction(
			ctx,
			data[0],
			data[1:len(data)-1],
		)

		if err != nil {
			s.logger.LogAttrs(
				ctx,
				slog.LevelError,
				"action error",
				slog.Any("error", err),
			)
		}
	}
}

func (s *serverImpl) processAction(
	ctx context.Context,
	actionType byte,
	data []byte,
) error {
	switch actionType {
	case models.GoMaxProcOperation:
		action := models.GoMaxProcAction{}
		err := json.Unmarshal(data, &action)

		if err != nil {
			return err
		}

		s.goMaxProcOp(ctx, &action)
	case models.SetMemoryLimitOperation:
		action := models.SetMemoryLimitAction{}
		err := json.Unmarshal(data, &action)

		if err != nil {
			return err
		}

		s.setMemoryLimitOp(ctx, &action)
	case models.DoRequestsOperation:
		action := models.DoRequestAction{}
		err := json.Unmarshal(data, &action)

		if err != nil {
			return err
		}

		s.doRequestOp(ctx, &action)
	default:
		return fmt.Errorf("unsupported operation %b", actionType)
	}

	return nil
}

func (s *serverImpl) invariantError(ctx context.Context, err error) {
	s.logger.LogAttrs(
		ctx,
		slog.LevelError,
		"invariant error",
		slog.Any("error", err),
	)
}

func (s *serverImpl) sendResultMessage(ctx context.Context, actionType byte, result any) {
	data, err := json.Marshal(result)

	if err != nil {
		s.invariantError(ctx, err)
		return
	}

	err = s.output.WriteByte(actionType)

	if err != nil {
		s.invariantError(ctx, err)
		return
	}

	_, err = s.output.Write(data)

	if err != nil {
		s.invariantError(ctx, err)
		return
	}

	err = s.output.WriteByte(models.ActionsDelimiter)

	if err != nil {
		s.invariantError(ctx, err)
		return
	}

	err = s.output.Flush()

	if err != nil {
		s.invariantError(ctx, err)
		return
	}
}
