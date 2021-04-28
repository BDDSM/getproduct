package productProvider

import (
	"github.com/korableg/getproduct/pkg/product"
)

type ProductProvider interface {
	GetProduct(cbarcode string) (*product.Product, error)
}
