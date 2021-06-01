package localprovider

import (
	"context"

	"github.com/korableg/getproduct/pkg/product"
	productProvider "github.com/korableg/getproduct/pkg/product/provider"
)

type ProductLocalProvider interface {
	productProvider.ProductProvider
	AddProduct(ctx context.Context, p *product.Product) error
	DeleteProduct(ctx context.Context, barcode string) error
}
