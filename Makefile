#   Copyright 2024 Pranav Singh

#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at

#   http://www.apache.org/licenses/LICENSE-2.0

#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

include install/Makefile.core.mk
include install/Makefile.help.mk

## ---------------------------
## LOCAL BUILD
## ---------------------------
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


## ---------------------------
## DOCKER BUILD
## ---------------------------
.PHONY: docker-build

## Builds the docker image
docker-build:
	docker buildx build --platform=linux/arm64 -f ./install/docker/Dockerfile -t $(BINARY_NAME):latest .