.PHONY: build test clean

# The name of the binary to produce
BINARY_NAME=hippocampus

# Takes care of all the dependencies
deps:
	go get -v -d ./...

# Runs tests
test: deps
	go test -v ./...

# Builds the binary
build:
	cd cmd; \
  go build -o $(BINARY_NAME)

# Runs the server without building the binary
server:
	cd cmd; \
	go run main.go

# Runs the binary
run:
	cd cmd; \
	./$(BINARY_NAME)

# Cleans up the binary
clean:
	go clean; \
	rm -f $(BINARY_NAME)

# Default make will run all
all: deps test build