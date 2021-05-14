package productLocalProvider

import (
	"context"
	"github.com/korableg/getproduct/pkg/product"
	"github.com/korableg/getproduct/pkg/productProvider"
)

type ProductLocalProvider interface {
	productProvider.ProductProvider
	AddProduct(ctx context.Context, product *product.Product) error
	DeleteProduct(ctx context.Context, product *product.Product) error
}