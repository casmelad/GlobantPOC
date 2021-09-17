package main

import (
	"fmt"
	"log"
	"net"

	"github.com/caarlos0/env/v6"
	uservice "github.com/casmelad/GlobantPOC/cmd/grpc_server/users"
	"google.golang.org/grpc"
)

func main() {

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))

	if err != nil {
		log.Fatalf("Could not create the listener %v", err)
	}

	server := grpc.NewServer()
	uservice.RegisterUsersServer(server, uservice.NewUserService())

	if err := server.Serve(ls); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	fmt.Println("Server is running!")
}

type config struct {
	Port int `env:"GRPCSERVICE_PORT" envDefault:"9000"`
}
