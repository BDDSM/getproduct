// Package classification of GetProduct API
//
// Documentanion for Product API
//
//  Schemes: http
//  BasePath: /api
//  Version: 1.0.1.4
//
//  Produces:
//  - application/json
// swagger:meta
package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korableg/getproduct/pkg/errs"
	lp "github.com/korableg/getproduct/pkg/product/localprovider"
	pp "github.com/korableg/getproduct/pkg/product/provider"
	productRepository "github.com/korableg/getproduct/pkg/product/repository"
)

var engine *gin.Engine
var repository *productRepository.ProductRepository

// swagger:response noContent
type noContent struct {
}

func init() {

	repository = productRepository.New(nil)
	// repository.AddProvider(&barcodeList.BarcodeList{})

	// repository.AddProvider(&vekaptek.Vekaptek{})
	// repository.AddProvider(&disai.Disai{})

	// if config.ChromeDPConfig() != nil {

	// 	chromeDPWSAddress := fmt.Sprintf("ws://%s:%d", config.ChromeDPConfig().Hostname, config.ChromeDPConfig().Port)

	// 	repository.AddProvider(nationalCatalog.New(chromeDPWSAddress))
	// 	repository.AddProvider(biostyle.New(chromeDPWSAddress))
	// 	//repository.AddProvider(&eapteka.Eapteka{})
	// }

	engine = gin.New()
	engine.Use(defaultHeaders)

	engine.NoRoute(pageNotFound)
	engine.NoMethod(methodNotAllowed)

	group := engine.Group("/api/barcode")
	group.Use(validateRequest)

	group.GET("/first/:barcode", getProduct)
	group.GET("/thebest/:barcode", getTheBestProduct)
	group.GET("/all/:barcode", getAllProducts)

	engine.DELETE("/api/localstorage/:barcode", deleteProductFromLocalRepository)

}

func Engine() *gin.Engine {
	return engine
}

func SetLocalProvider(lp lp.ProductLocalProvider) {
	repository.SetLocalProvider(lp)
}

func AddProvider(provider pp.ProductProvider) {
	repository.AddProvider(provider)
}

func defaultHeaders(c *gin.Context) {
	c.Next()
	c.Header("Server", fmt.Sprintf("GetProduct:%s", "1.0.1.4"))
}

func validateRequest(c *gin.Context) {
	barcode := c.Params.ByName("barcode")
	if barcode == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, errs.New(errors.New("barcode hasn't filled")))
	}
}

func pageNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, errs.New(errors.New("not found")))
}

func methodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, errs.New(errors.New("method is not allowed")))
}

// swagger:route GET /api/barcode/first/:barcode product firstProduct
// Returns first found product by barcode
// responses:
//  200: product
//  400: error
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

// swagger:route GET /api/barcode/thebest/:barcode product theBestProduct
// Returns the best found product by barcode
// responses:
//  200: product
//  400: error
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

// swagger:route GET /api/barcode/all/:barcode product allProduct
// Returns all found variants of product by barcode
// responses:
//  200: []product
//  400: error
func getAllProducts(c *gin.Context) {
	barcode := c.Params.ByName("barcode")

	p, err := repository.GetAll(c, barcode)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, errs.New(err))
		return
	}

	c.JSON(http.StatusOK, p)
}

// swagger:route DELETE /api/localstorage/:barcode product delete
// Deletes product by barcode from local storage
// responses:
//  200: noContent
//  400: error
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
