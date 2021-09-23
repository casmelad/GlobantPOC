FROM alpine as grpcserver

ENV GO111MODULE=on

RUN apk add --no-cache go
RUN apk add bash ca-certificates git gcc g++ libc-dev

WORKDIR /go/src/grpc
COPY . .

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go build -o grpcservice ./cmd/grpc_server/.

EXPOSE 9000

CMD ./grpcservice


# Set necessary environmet variables needed for the REST server

FROM alpine as restserver


ENV GO111MODULE=on

RUN apk add --no-cache go
RUN apk add bash ca-certificates git gcc g++ libc-dev

WORKDIR /go/src/rest
COPY . .

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go build -o restservice ./cmd/REST_server/.

EXPOSE 8000

CMD ./restservice