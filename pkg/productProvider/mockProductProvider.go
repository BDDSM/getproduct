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
		return product.New("111", "http://testurl.ru"), nil
	}

	return nil, errors.New("product didn't find by barcode")

}
