include install/Makefile.core.mk
include install/Makefile.help.mk

.PHONY: build test clean

## Takes care of all the dependencies
deps:
	go get -v -d ./...

## Runs tests
test: deps
	go test -v ./...

## Builds the binary
build:
	cd cmd; go mod tidy; \
  go build -o $(BINARY_NAME)

## Runs the server without building the binary
server:
	cd cmd; go mod tidy; \
	go run main.go

## Runs the binary
run:
	cd cmd; \
	./$(BINARY_NAME)

## Cleans up the binary
clean:
	go clean; \
	rm -f $(BINARY_NAME)

## Default make will run all
all: deps test build