package main

import (
	"fmt"
	"log"
	"net"

	uservice "github.com/casmelad/GlobantPOC/cmd/grpc_server/users"
	"google.golang.org/grpc"
)

func main() {

	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))

	if err != nil {
		log.Fatalf("Could not create the listener %v", err)
	}

	server := grpc.NewServer()

	uservice.RegisterUsersServer(server, uservice.NewUserService())

	if err := server.Serve(ls); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	fmt.Println("Server is up!")
}
