OS := $(shell uname)

run: build
	./bin/${WATCHDOG_BINARY_NAME}

build: init-binary-names
	go mod tidy
	go build -o ./bin/${SERVER_BINARY_NAME} ./cmd/server/server.go
	go build -o ./bin/${WATCHDOG_BINARY_NAME} ./cmd/watchdog/watchdog.go

init-binary-names:
ifeq (${OS},Windows_NT)
SERVER_BINARY_NAME=server.exe
WATCHDOG_BINARY_NAME=watchdog.exe
else
SERVER_BINARY_NAME=server
WATCHDOG_BINARY_NAME=watchdog
endif
