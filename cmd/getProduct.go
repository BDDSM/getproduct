package main

import (
	"fmt"

	httpServer "github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/config"
	"github.com/asim/go-micro/v3/config/source/env"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/korableg/getproduct/internal/api"
)

const nameApp = "getproduct"

func init() {
	godotenv.Load()
	config.Load(
		env.NewSource(
			env.WithStrippedPrefix("getproduct"),
		),
	)
}

func main() {

	if config.Get("debug").Bool(false) {
		logger.Init(logger.WithLevel(logger.DebugLevel))
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	address := fmt.Sprintf("%s:%d", config.Get("address").String(""), config.Get("port").Int(11218))

	srv := httpServer.NewServer(
		server.Name(nameApp),
		server.Address(address),
	)

	hd := srv.NewHandler(api.Engine())
	if err := srv.Handle(hd); err != nil {
		logger.Fatal(err)
	}

	service := micro.NewService(
		micro.Name(nameApp),
		micro.Server(srv),
		micro.Registry(registry.NewRegistry()),
	)

	service.Init()
	service.Run()

}
