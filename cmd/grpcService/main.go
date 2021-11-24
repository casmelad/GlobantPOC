package main

import (
	"fmt"
	"net"
	"os"

	"github.com/go-kit/kit/log"

	"github.com/caarlos0/env/v6"
	grpcServiceImpl "github.com/casmelad/GlobantPOC/cmd/grpcService/users"
	proto "github.com/casmelad/GlobantPOC/cmd/grpcService/users/proto"
	memory "github.com/casmelad/GlobantPOC/pkg/repository/memory"
	mysql "github.com/casmelad/GlobantPOC/pkg/repository/mysql"
	domain "github.com/casmelad/GlobantPOC/pkg/users"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"
)

func main() {

	var (
		zipkinURL = "http://localhost:9411/api/v2/spans"
	)

	// Create a single logger, which we'll use and give to other components.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var zipkinTracer *zipkin.Tracer
	{
		if zipkinURL != "" {
			var (
				err         error
				hostPort    = "localhost:8081"
				serviceName = "accounts"
				reporter    = zipkinhttp.NewReporter(zipkinURL)
			)
			defer reporter.Close()
			//sampler, err := zipkin.NewCountingSampler(100)
			zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
			zipkinTracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zEP))
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
		}
	}

	// Determine which OpenTracing tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency.
	var tracer stdopentracing.Tracer
	/* {

		logger.Log("tracer", "Zipkin", "type", "OpenTracing", "URL", zipkinURL)
		tracer = zipkinot.Wrap(zipkinTracer)
		zipkinTracer = nil // do not instrument with both native tracer and opentracing bridge

	} */

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))

	if err != nil {
		panic(fmt.Sprintf("Could not create the listener %v", err))
	}

	userService := domain.NewUserService(getActiveRepository())
	endpoints := grpcServiceImpl.NewGrpcUsersServer(userService)

	grpcUserServer := grpcServiceImpl.NewGrpcUserServer(*endpoints, tracer, zipkinTracer, logger)

	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	proto.RegisterUsersServer(baseServer, grpcUserServer)

	if err := baseServer.Serve(ls); err != nil {
		panic(fmt.Sprintf("failed to serve: %s", err))
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
			panic(fmt.Sprintf("mysql connection failed: %s", err))
		}
		return repo
	}
	return nil
}

type config struct {
	Port int `env:"GRPCSERVICE_PORT" envDefault:"9000"`
}
