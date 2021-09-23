package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/casmelad/GlobantPOC/cmd/REST_server/web"
	"github.com/gorilla/mux"
)

func main() {

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	router := mux.NewRouter()

	fmt.Println(cfg)

	web.Startup(router)

	srv := &http.Server{
		Handler: router,
		Addr:    cfg.Hosts + strconv.Itoa(cfg.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
	}

	err := srv.ListenAndServe()

	fmt.Println(err)

	if err != nil {
		log.Fatal(err)
	}

}

type config struct {
	Port         int    `env:"RESTSERVER_PORT" envDefault:"8000"`
	Hosts        string `env:"RESTSERVER_HOSTS" envDefault:"127.0.0.1:"`
	WriteTimeout int    `env:"RESTSERVER_WRITETIMEOUT" envDefault:"15"`
	ReadTimeout  int    `env:"RESTSERVER_READTIMEOUT" envDefault:"15"`
}
