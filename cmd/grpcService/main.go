package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/caarlos0/env/v6"
	grpcServiceImpl "github.com/casmelad/GlobantPOC/cmd/grpcService/users"
	proto "github.com/casmelad/GlobantPOC/cmd/grpcService/users/proto"
	memory "github.com/casmelad/GlobantPOC/pkg/repository/memory"
	mysql "github.com/casmelad/GlobantPOC/pkg/repository/mysql"
	domain "github.com/casmelad/GlobantPOC/pkg/users"
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
	appService := domain.NewUserService(getActiveRepository())
	grpcService := grpcServiceImpl.NewGrpcUserService(appService)
	proto.RegisterUsersServer(server, grpcService)

	if err := server.Serve(ls); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func getActiveRepository() domain.Repository {

	envVar := os.Getenv("USERS_REPOSITORY")

	fmt.Println(envVar)

	if len(envVar) == 0 {
		envVar = "mysql"
	}

	switch envVar {
	case "memory":
		repo := memory.NewInMemoryUserRepository()
		return repo
	case "mysql":
		repo, err := mysql.NewMySQLUserRepository()
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}

type config struct {
	Port int `env:"GRPCSERVICE_PORT" envDefault:"9000"`
}
