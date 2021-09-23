FROM alpine as grpcserver

RUN apk add --no-cache go
RUN go version
WORKDIR /go/src/grpc
COPY . .

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go build -o grpcservice ./cmd/grpc_server/.

EXPOSE 9000

CMD ["/go/src/grpc/grpcservice"]


# Set necessary environmet variables needed for the REST server

FROM alpine as restserver


RUN apk add --no-cache go
RUN go version

WORKDIR /go/src/rest
COPY . .

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go build -o restservice ./cmd/REST_server/.

EXPOSE 8000

CMD ["/go/src/rest/restservice"]