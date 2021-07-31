package repository

import (
	"reflect"

	"github.com/korableg/getproduct/pkg/product/localprovider"
	"github.com/korableg/getproduct/pkg/product/provider"
)

type ProductRepositoryOption func(*ProductRepository)

func WithLocalProvider(lp localprovider.ProductLocalProvider) ProductRepositoryOption {
	return func(p *ProductRepository) {
		if !reflect.ValueOf(lp).IsNil() {
			p.localProvider = lp
		}
	}
}

func WithChromeDP(chromeDPWSAddress string) ProductRepositoryOption {
	return func(p *ProductRepository) {
		p.chromeDPWSAddress = chromeDPWSAddress
	}
}

func WithProviders(pr ...provider.ProductProvider) ProductRepositoryOption {
	return func(p *ProductRepository) {
		if !reflect.ValueOf(pr).IsNil() {
			p.providers = append(p.providers, pr...)
		}
	}
}
