package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korableg/getproduct/internal/config"
	"github.com/korableg/getproduct/internal/errs"
	"github.com/korableg/getproduct/pkg/localProviders/mongod"
	"github.com/korableg/getproduct/pkg/productLocalProvider"
	"github.com/korableg/getproduct/pkg/productProviders/barcodeList"
	"github.com/korableg/getproduct/pkg/productProviders/biostyle"
	"github.com/korableg/getproduct/pkg/productProviders/disai"
	"github.com/korableg/getproduct/pkg/productProviders/nationalCatalog"
	"github.com/korableg/getproduct/pkg/productProviders/vekaptek"
	"github.com/korableg/getproduct/pkg/productRepository"
)

var engine *gin.Engine
var repository *productRepository.ProductRepository

func init() {

	var localProvider productLocalProvider.ProductLocalProvider

	mongoCfg := config.MongoDBConfig()
	if mongoCfg != nil {
		var err error
		if localProvider, err = mongod.NewMongoDB(mongoCfg.Hostname, mongoCfg.Port, mongoCfg.Username, mongoCfg.Password); err != nil {
			log.Println(err)
			panic(err)
		}
	}

	repository = productRepository.New(localProvider)
	repository.AddProvider(&barcodeList.BarcodeList{})

	repository.AddProvider(&vekaptek.Vekaptek{})
	repository.AddProvider(&disai.Disai{})

	if config.ChromeDPConfig() != nil {

		chromeDPWSAddress := fmt.Sprintf("ws://%s:%d", config.ChromeDPConfig().Hostname, config.ChromeDPConfig().Port)

		repository.AddProvider(nationalCatalog.New(chromeDPWSAddress))
		repository.AddProvider(biostyle.New(chromeDPWSAddress))
		//repository.AddProvider(&eapteka.Eapteka{})
	}

	if config.Debug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine = gin.New()
	engine.Use(defaultHeaders())

	engine.NoRoute(pageNotFound)
	engine.NoMethod(methodNotAllowed)

	engine.GET("/api/barcode/:barcode", getProduct)
	engine.GET("/api/thebestproduct/:barcode", getTheBestProduct)

	engine.DELETE("/api/localstorage/:barcode", deleteProductFromLocalRepository)

}

func Run() {
	go func() {
		address := fmt.Sprintf("%s:%d", config.Address(), config.Port())
		err := engine.Run(address)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func defaultHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Server", fmt.Sprintf("GetProduct:%s", config.Version()))
	}
}

func pageNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, errs.New(errors.New("not found")))
}

func methodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, errs.New(errors.New("method is not allowed")))
}

func getProduct(c *gin.Context) {
	barcode := c.Params.ByName("barcode")

	p, err := repository.Get(c, barcode)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, errs.New(err))
		return
	}

	c.JSON(http.StatusOK, p)
}

func getTheBestProduct(c *gin.Context) {
	barcode := c.Params.ByName("barcode")

	p, err := repository.GetTheBest(c, barcode)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, errs.New(err))
		return
	}

	c.JSON(http.StatusOK, p)
}

func deleteProductFromLocalRepository(c *gin.Context) {
	barcode := c.Params.ByName("barcode")

	err := repository.DeleteFromLocalProvider(c, barcode)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, errs.New(err))
		return
	}

	c.Status(http.StatusOK)
}
