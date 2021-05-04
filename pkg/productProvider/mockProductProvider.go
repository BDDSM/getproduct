package productProvider

import (
	"context"
	"errors"

	"github.com/korableg/getproduct/pkg/product"
)

type mockProductProvider struct {
}

func (m *mockProductProvider) GetProduct(ctx context.Context, barcode string) (*product.Product, error) {

	if barcode == "111" {
		return product.NewProduct("111", "TestProduct", "шт", "TestDescription", "TestM"), nil
	}

	return nil, errors.New("product didn't find by barcode")

}
