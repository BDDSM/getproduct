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
	productRepository "github.com/korableg/getproduct/pkg/product/repository"
)

// swagger:response noContent
type noContent struct {
}

type Engine struct {
	engine     *gin.Engine
	repository *productRepository.ProductRepository
}

func New(opts ...EngineOption) (*Engine, error) {

	e := &Engine{}

	engine := gin.New()
	engine.Use(e.defaultHeaders)

	engine.NoRoute(e.pageNotFound)
	engine.NoMethod(e.methodNotAllowed)

	group := engine.Group("/api/barcode")
	group.Use(e.validateRequest)

	group.GET("/first/:barcode", e.getProduct)
	group.GET("/thebest/:barcode", e.getTheBestProduct)
	group.GET("/all/:barcode", e.getAllProducts)

	engine.DELETE("/api/localstorage/:barcode", e.deleteProductFromLocalRepository)

	e.engine = engine

	for _, opt := range opts {
		opt(e)
	}

	if e.repository == nil {
		return nil, errors.New("product repository didn't initialize")
	}

	return e, nil

}

func init() {

	// repository.AddProvider(&barcodeList.BarcodeList{})

	// repository.AddProvider(&vekaptek.Vekaptek{})
	// repository.AddProvider(&disai.Disai{})

	// if config.ChromeDPConfig() != nil {

	// 	repository.AddProvider(nationalCatalog.New(chromeDPWSAddress))
	// 	repository.AddProvider(biostyle.New(chromeDPWSAddress))
	// 	//repository.AddProvider(&eapteka.Eapteka{})
	// }

}

func (e *Engine) Handler() *gin.Engine {
	return e.engine
}

func (e *Engine) defaultHeaders(c *gin.Context) {
	c.Next()
	c.Header("Server", fmt.Sprintf("GetProduct:%s", "1.0.2.1"))
}

func (e *Engine) validateRequest(c *gin.Context) {
	barcode := c.Params.ByName("barcode")
	if barcode == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, errs.New(errors.New("barcode hasn't filled")))
	}
}

func (e *Engine) pageNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, errs.New(errors.New("not found")))
}

func (e *Engine) methodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, errs.New(errors.New("method is not allowed")))
}

// swagger:route GET /api/barcode/first/:barcode product firstProduct
// Returns first found product by barcode
// responses:
//  200: product
//  400: error
func (e *Engine) getProduct(c *gin.Context) {
	barcode := c.Params.ByName("barcode")

	p, err := e.repository.Get(c, barcode)
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
func (e *Engine) getTheBestProduct(c *gin.Context) {
	barcode := c.Params.ByName("barcode")

	p, err := e.repository.GetTheBest(c, barcode)
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
func (e *Engine) getAllProducts(c *gin.Context) {
	barcode := c.Params.ByName("barcode")

	p, err := e.repository.GetAll(c, barcode)
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
func (e *Engine) deleteProductFromLocalRepository(c *gin.Context) {
	barcode := c.Params.ByName("barcode")

	err := e.repository.DeleteFromLocalProvider(c, barcode)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, errs.New(err))
		return
	}

	c.Status(http.StatusOK)
}
