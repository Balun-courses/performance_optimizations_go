package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"hw1/models"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"unsafe"
)

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

type watchdog struct {
	cliReader    reader
	serverReader reader
	serverWriter writer
	logger       *slog.Logger
}

func newWatchDog(
	cliReader reader,
	serverReader reader,
	serverWriter writer,
	logger *slog.Logger,
) *watchdog {
	return &watchdog{
		cliReader:    cliReader,
		serverReader: serverReader,
		serverWriter: serverWriter,
		logger:       logger,
	}
}

// TODO: add more log information for students

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rd := bufio.NewReader(os.Stdin)

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	sp := flag.String("sp", "", "path to server binary")
	serverArguments := flag.String("sargs", "", "server program arguments")
	flag.Parse()

	if *sp == "" {
		logger.Error("server_path -sp must be set for watchdog actions")
		return
	}

	if *serverArguments == "" {
		logger.Error("server_arguments -sargs must be set for watchdog actions")
	}

	serverReader, serverWriter, err := os.Pipe()

	if err != nil {
		logger.Error("can not create pipe", slog.Any("error", err))
		return
	}

	go func() {
		<-ctx.Done()

		err = serverReader.Close()

		if err != nil {
			panic(err)
		}
	}()

	watchdogReader, watchdogWriter, err := os.Pipe()

	if err != nil {
		logger.Error("can not create pipe", slog.Any("error", err))
		return
	}

	watchDog := newWatchDog(
		rd,
		bufio.NewReaderSize(watchdogReader, models.MaxActionSize),
		bufio.NewWriter(serverWriter),
		logger,
	)

	cmd := exec.CommandContext(ctx, *sp)
	cmd.Stdin = serverReader
	cmd.Stdout = watchdogWriter
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Args = append(cmd.Args, strings.Split(*serverArguments, " ")...)

	err = cmd.Start()

	if err != nil {
		logger.Error("can not start server", slog.Any("error", err))
		return
	}

	fmt.Println("Server start successfully")

	for {
		var commandType string

		_, err = fmt.Fscan(rd, &commandType)

		if err != nil {
			logger.Error("can not fetch input command", slog.Any("error", err))
			break
		}

		var cErr error

		switch commandType {
		case models.WatchdogCommandDoRequest:
			cErr = watchDog.processDoReq()
		case models.WatchDogCommandSetMacProc:
			cErr = watchDog.processSetMaxProc()
		case models.WatchDogCommandCommandSetMemBytes:
			cErr = watchDog.processSetMemLimit()
		case models.WatchDogCommandExit:
			return
		default:
			logger.Error("unsupported command", slog.String("command", commandType))
			continue
		}

		if cErr != nil {
			logger.Error("an error occurred while executing the command",
				slog.String("command", commandType),
				slog.Any("error", err),
			)
			break
		}
	}

	if err != nil {
		logger.Error("can not close server serverWriter", slog.Any("error", err))
		return

	}

	err = cmd.Wait()

	if err != nil {
		panic(err)
	}
}

func (w *watchdog) processDoReq() error {
	actionLogger := w.logger.With(slog.String("action", "processDoReq"))

	var (
		method string
		sUrl   string
		body   string
	)

	_, err := fmt.Fscan(w.cliReader, &method, &sUrl, &body)

	if err != nil {
		actionLogger.Error("can not read action params",
			slog.Any("error", err),
		)

		return err
	}

	serverUrl, err := url.Parse(sUrl)

	if err != nil {
		actionLogger.Error("can not parse server url", slog.Any("error", err))
		return err
	}

	action := models.DoRequestAction{
		RequestID: rand.N[int64](123_000),
		Url:       serverUrl,
		Method:    method,
		Body:      unsafe.Slice(unsafe.StringData(body), len(body)),
	}

	var result models.DoRequestActionResult

	err = doAction(
		w.serverReader,
		w.serverWriter,
		w.logger,
		models.DoRequestsOperation,
		action,
		&result,
	)

	fmt.Printf("processDoReq result: %v\n", &result)

	return nil
}

func (w *watchdog) processSetMaxProc() error {
	actionLogger := w.logger.With(slog.String("action", "processSetMaxProc"))

	var value int

	_, err := fmt.Fscan(w.cliReader, &value)

	if err != nil {
		actionLogger.Error("can not read action params",
			slog.Any("error", err),
		)

		return err
	}

	action := models.GoMaxProcAction{
		Value: value,
	}

	var result models.GoMaxProcActionResult

	err = doAction(
		w.serverReader,
		w.serverWriter,
		w.logger,
		models.GoMaxProcOperation,
		action,
		&result,
	)

	fmt.Printf("processSetMaxProc result: %v\n", &result)

	return nil
}

func (w *watchdog) processSetMemLimit() error {
	actionLogger := w.logger.With(slog.String("action", "processSetMemLimit"))

	var value int64

	_, err := fmt.Fscan(w.cliReader, &value)

	if err != nil {
		actionLogger.Error("can not read action params",
			slog.Any("error", err),
		)

		return err
	}

	action := models.SetMemoryLimitAction{
		Value: value,
	}

	var result models.SetMemoryLimitActionResult

	err = doAction(
		w.serverReader,
		w.serverWriter,
		w.logger,
		models.SetMemoryLimitOperation,
		action,
		&result,
	)

	fmt.Printf("processSetMemLimit result: %v\n", &result)

	return nil
}

func doAction(
	serverReader reader,
	serverWriter writer,
	actionLogger *slog.Logger,
	actionType byte,
	action any,
	response any,
) error {
	serialized, err := json.Marshal(action)

	if err != nil {
		actionLogger.Error("can not marshal action", slog.Any("error", err))
		return err
	}

	err = serverWriter.WriteByte(actionType)

	if err != nil {
		actionLogger.Error("can not write to server", slog.Any("error", err))
		return err
	}

	_, err = serverWriter.Write(serialized)

	if err != nil {
		actionLogger.Error("can not write to server", slog.Any("error", err))
		return err
	}

	err = serverWriter.WriteByte(models.ActionsDelimiter)

	if err != nil {
		actionLogger.Error("can not write to server", slog.Any("error", err))
		return err
	}

	err = serverWriter.Flush()

	if err != nil {
		actionLogger.Error("can not flush server writer", slog.Any("error", err))
		return err
	}

	// sync here, may be better
	res, err := serverReader.ReadSlice(models.ActionsDelimiter)

	if err != nil {
		actionLogger.Error("can reader server result", slog.Any("error", err))
		return err
	}

	if res[0] != actionType {
		actionLogger.Error("invariant error", slog.Any("error", err))
		return fmt.Errorf("invariant error")
	}

	err = json.Unmarshal(res[1:len(res)-1], &response)

	if err != nil {
		actionLogger.Error("can not unmarshal action result", slog.Any("error", err))
		return err
	}

	return nil
}
