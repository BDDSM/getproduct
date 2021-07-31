package provider

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/korableg/getproduct/pkg/product"
)

var ErrProductDidntFind error = errors.New("product didn't find")
var ChromeDPWSAddress string

var providers map[string]ProductProvider
var providersMu *sync.RWMutex

func init() {
	providers = make(map[string]ProductProvider)
	providersMu = &sync.RWMutex{}
}

type ProductProvider interface {
	GetProduct(ctx context.Context, barcode string) (*product.Product, error)
}

func Register(name string, p ProductProvider) {
	providersMu.Lock()
	defer providersMu.Unlock()

	if _, ok := providers[name]; ok {
		panic(errors.New("provider has been registered"))
	}

	providers[name] = p

}

func Get(name string) ProductProvider {

	providersMu.RLock()
	defer providersMu.RUnlock()

	p, ok := providers[name]

	if !ok {
		panic(fmt.Errorf("provider by name %s hasn't registered", name))
	}

	return p

}

func GetAll() []ProductProvider {

	providersMu.RLock()
	defer providersMu.RUnlock()

	allProviders := make([]ProductProvider, len(providers))

	i := 0
	for _, v := range providers {
		allProviders[i] = v
		i++
	}

	return allProviders

}
