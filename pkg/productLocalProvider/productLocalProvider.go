package productLocalProvider

import (
	"context"
	"github.com/korableg/getproduct/pkg/product"
	"github.com/korableg/getproduct/pkg/productProvider"
)

type ProductLocalProvider interface {
	productProvider.ProductProvider
	AddProduct(ctx context.Context, p *product.Product) error
	DeleteProduct(ctx context.Context, barcode string) error
}
