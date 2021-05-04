package productRepository

import (
	"context"
	"errors"
	"fmt"
	"github.com/korableg/getproduct/pkg/product"
	"github.com/korableg/getproduct/pkg/productProvider"
	"log"
	"sync"
	"time"
)

type ProductRepository struct {
	providers   []productProvider.ProductProvider
	muProviders sync.RWMutex
}

func NewProductRepository() *ProductRepository {
	pr := ProductRepository{
		providers: make([]productProvider.ProductProvider, 0, 10),
	}

	return &pr
}

func (pr *ProductRepository) AddProvider(provider productProvider.ProductProvider) {
	pr.muProviders.Lock()
	defer pr.muProviders.Unlock()
	pr.providers = append(pr.providers, provider)
}

func (pr *ProductRepository) Get(ctx context.Context, barcode string) (*product.Product, error) {
	pr.muProviders.RLock()
	defer pr.muProviders.RUnlock()

	log.Println(fmt.Sprintf("getting product by barcode: %s", barcode))

	if len(pr.providers) == 0 {
		return nil, errors.New("product providers is empty")
	}

	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer cancelFunc()

	productChan := make(chan *product.Product)

	for _, provider := range pr.providers {
		go func() {
			p, err := provider.GetProduct(newCtx, barcode)
			if err != nil {
				log.Println(err)
				return
			}
			productChan <- p
		}()
	}

	select {
	case dst := <-productChan:
		return dst, nil
	case <-newCtx.Done():
		return nil, newCtx.Err()
	}

}
