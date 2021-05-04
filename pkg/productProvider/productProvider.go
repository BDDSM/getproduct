package productProvider

import (
	"context"
	"errors"
	"github.com/korableg/getproduct/pkg/product"
)

var ErrProductDidntFind error = errors.New("product didn't find")

type ProductProvider interface {
	GetProduct(ctx context.Context, barcode string) (*product.Product, error)
}
