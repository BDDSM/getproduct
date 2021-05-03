package productProvider

import (
	"errors"
	"github.com/korableg/getproduct/pkg/product"
)

var ErrProductDidntFind error = errors.New("product didn't find")

type ProductProvider interface {
	GetProduct(cbarcode string) (*product.Product, error)
}
