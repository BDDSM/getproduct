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
	"github.com/korableg/getproduct/pkg/product/localprovider/mongod"
	"github.com/korableg/getproduct/pkg/product/provider"
	_ "github.com/korableg/getproduct/pkg/product/provider/barcodeList"
	_ "github.com/korableg/getproduct/pkg/product/provider/biostyle"
	_ "github.com/korableg/getproduct/pkg/product/provider/disai"
	_ "github.com/korableg/getproduct/pkg/product/provider/nationalCatalog"
	_ "github.com/korableg/getproduct/pkg/product/provider/vekaptek"
	"github.com/korableg/getproduct/pkg/product/repository"
	"github.com/pkg/errors"
)

const nameApp = "getproduct"

func init() {
	godotenv.Load()
}

func main() {

	var mongodb *mongod.MongoDB
	var chromeDPWSAddress string
	var err error

	src := env.NewSource(
		env.WithStrippedPrefix("GETPRODUCT"),
	)

	if err := config.Load(src); err != nil {
		logger.Fatal(errors.Wrap(err, "config loading"))
	}

	if config.Get("debug").Bool(false) {
		logger.Init(logger.WithLevel(logger.DebugLevel))
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if config.Get("mongodb", "use").Bool(false) {
		if mongodb, err = mongod.NewMongoDB(
			config.Get("mongodb", "hostname").String("localhost"),
			config.Get("mongodb", "port").Int(27017),
			config.Get("mongodb", "username").String(""),
			config.Get("mongodb", "password").String("")); err != nil {
			panic(err)
		}
	}

	if config.Get("chromedp", "use").Bool(false) {
		chromeDPWSAddress = fmt.Sprintf("ws://%s:%d",
			config.Get("chromedp", "hostname").String("localhost"), config.Get("chromedp", "port").Int(3000))
	}

	pr := repository.New(
		repository.WithLocalProvider(mongodb),
		repository.WithProviders(provider.GetAll()...),
		repository.WithChromeDP(chromeDPWSAddress),
	)

	engine, err := api.New(
		api.WithProductRepository(pr),
	)

	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf("%s:%d", config.Get("address").String(""), config.Get("port").Int(11218))

	srv := httpServer.NewServer(
		server.Name(nameApp),
		server.Address(address),
	)

	hd := srv.NewHandler(engine.Handler())
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
