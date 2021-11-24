############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR /src
COPY . .
# Fetch dependencies.
# Using go get.
RUN go mod download
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o grpcservice ./cmd/grpcService/.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o restservice ./cmd/restService/.

############################
# STEP 2 grpc service
############################
FROM golang:alpine as grpcserver

WORKDIR /

COPY --from=builder /src/grpcservice ./grpcservice

ENTRYPOINT ["/grpcservice"]

############################
# STEP 3 REST service
############################
FROM golang:alpine as restserver

WORKDIR /

COPY --from=builder /src/restservice ./restservice

ENTRYPOINT ["/restservice"]