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

# syntax=docker/dockerfile:1

FROM golang:1.21 as builder

RUN mkdir -p /app
WORKDIR /app

COPY . .
RUN cd cmd && GOSUMDB=off GOARCH=arm64 CGO_ENABLED=0 go build -o ./hippocampus .

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/cmd/hippocampus /app/hippocampus

# Create data directory
RUN mkdir -p /data

# Default environment variables
ENV HIPPOCAMPUS_PORT=6379
ENV HIPPOCAMPUS_MODE=single
ENV HIPPOCAMPUS_NODE_ID=node-1
ENV HIPPOCAMPUS_DATA_DIR=/data
ENV RAFT_ENABLED=true
ENV RAFT_RPC_PORT=7000
ENV STORAGE_AOF_FILE=/data/database.aof
ENV STORAGE_ENABLE_AOF=true

EXPOSE 6379 7000

CMD ["/app/hippocampus"]


