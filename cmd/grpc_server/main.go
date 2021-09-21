package main

import (
	"fmt"
	"log"
	"net"

	"github.com/caarlos0/env/v6"
	grpcServiceImpl "github.com/casmelad/GlobantPOC/cmd/grpc_server/grpcservices"
	grpcServices "github.com/casmelad/GlobantPOC/cmd/grpc_server/users"
	infrastructure "github.com/casmelad/GlobantPOC/pkg/infrastructure"
	appservices "github.com/casmelad/GlobantPOC/pkg/services"
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
	repo := infrastructure.NewInMemoryUserRepository()
	appService := appservices.NewUserService(repo)
	grpcServices.RegisterUsersServer(server, grpcServiceImpl.NewGrpcUserService(appService))

	if err := server.Serve(ls); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	fmt.Println("Server is running!")
}

type config struct {
	Port int `env:"GRPCSERVICE_PORT" envDefault:"9000"` //Como pasarlo hacia abajo?
}
