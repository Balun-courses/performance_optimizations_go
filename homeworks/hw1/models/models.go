package models

import (
	"net/url"
)

const (
	WatchdogCommandDoRequest          = "doreq"
	WatchDogCommandSetMacProc         = "setproc"
	WatchDogCommandCommandSetMemBytes = "setmeml"
	WatchDogCommandExit               = "exit"
)

const (
	ServerBinaryPathArgName       = "sp"
	ServerProgramArgumentsArgName = "sargs"
)

const (
	ActionsDelimiter = '#' // admit a collision
	MaxActionSize    = 0xfff + 1<<10
)

const GoMaxProcOperation byte = 0x1

type (
	GoMaxProcAction struct {
		Value int
	}

	GoMaxProcActionResult struct {
		PreviousValue int
	}
)

const SetMemoryLimitOperation byte = 0x2

type (
	SetMemoryLimitAction struct {
		Value int64
	}

	SetMemoryLimitActionResult struct {
		PreviousValue int64
	}
)

const DoRequestsOperation byte = 0x3

type (
	DoRequestAction struct {
		RequestID int64
		Url       *url.URL
		Method    string
		Body      []byte
	}

	DoRequestActionResult struct {
		RequestID int64
		Error     string
		Response  []byte
	}
)
