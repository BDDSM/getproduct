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

	err := pr.checkProviders()
	if err != nil {
		return nil, err
	}

	log.Println(fmt.Sprintf("getting product by barcode: %s", barcode))

	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer cancelFunc()

	productChan := make(chan *product.Product)
	fetchingDoneChan := make(chan struct{})

	pr.getProductWithProviders(newCtx, barcode, productChan, fetchingDoneChan)

	select {
	case dst := <-productChan:
		return dst, nil
	case <-fetchingDoneChan:
		return nil, fmt.Errorf("product by barcode %s not found", barcode)
	case <-newCtx.Done():
		return nil, newCtx.Err()
	}

}

func (pr *ProductRepository) GetTheBest(ctx context.Context, barcode string) (*product.Product, error) {

	pr.muProviders.RLock()
	defer pr.muProviders.RUnlock()

	err := pr.checkProviders()
	if err != nil {
		return nil, err
	}

	log.Println(fmt.Sprintf("getting the best matches product by barcode: %s", barcode))

	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer cancelFunc()

	products := make([]*product.Product, 0, len(pr.providers))
	productChan := make(chan *product.Product)

	fetchingDoneChan := make(chan struct{})

	pr.getProductWithProviders(newCtx, barcode, productChan, fetchingDoneChan)

	for {
		select {
		case dst, ok := <-productChan:
			if ok {
				products = append(products, dst)
			} else {
				if len(products) > 0 {
					return pr.chooseTheBestProduct(products), nil
				} else {
					return nil, fmt.Errorf("product by barcode %s not found", barcode)
				}
			}

		case <-newCtx.Done():
			if len(products) > 0 {
				return pr.chooseTheBestProduct(products), nil
			}
			return nil, newCtx.Err()
		}
	}

}

func (pr *ProductRepository) getProductWithProviders(
	ctx context.Context, barcode string, productChan chan<- *product.Product, fetchingDoneChan chan<- struct{}) {

	wg := &sync.WaitGroup{}
	wg.Add(len(pr.providers))

	for _, provider := range pr.providers {
		go func(provider productProvider.ProductProvider, wg *sync.WaitGroup) {
			p, err := provider.GetProduct(ctx, barcode)
			defer wg.Done()
			if err != nil {
				log.Println(err)
				return
			}
			if p != nil {
				productChan <- p
			}

		}(provider, wg)
	}

	go func() {
		wg.Wait()
		close(productChan)
		fetchingDoneChan <- struct{}{}
	}()

}

func (pr *ProductRepository) checkProviders() error {
	if len(pr.providers) == 0 {
		return errors.New("product providers is empty")
	}

	return nil
}

func (pr *ProductRepository) chooseTheBestProduct(products []*product.Product) *product.Product {

	winner := products[0]
	winnerScore := winner.Rating()

	for i := 1; i < len(products); i++ {
		winnerCandidateRating := products[i].Rating()
		if winnerCandidateRating > winnerScore {
			winner = products[i]
			winnerScore = winnerCandidateRating
		}
	}

	return winner

}
