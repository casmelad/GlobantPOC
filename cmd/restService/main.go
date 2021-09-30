package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	glog "google.golang.org/grpc/grpclog"
)

var grpcLog glog.LoggerV2

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

func main() {

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	router := mux.NewRouter()

	fmt.Println(cfg)

	Startup(router)

	srv := &http.Server{
		Handler: router,
		Addr:    cfg.Hosts + strconv.Itoa(cfg.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
	}

	fmt.Println(srv.ListenAndServe())

	grpcLog.Info("Starting server at port :8000")
}

type config struct {
	Port         int    `env:"RESTSERVER_PORT" envDefault:"8000"`
	Hosts        string `env:"RESTSERVER_HOSTS" envDefault:":"`
	WriteTimeout int    `env:"RESTSERVER_WRITETIMEOUT" envDefault:"15"`
	ReadTimeout  int    `env:"RESTSERVER_READTIMEOUT" envDefault:"15"`
}
