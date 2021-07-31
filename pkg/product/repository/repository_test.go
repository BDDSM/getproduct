package repository

import (
	"context"
	"testing"

	_ "github.com/korableg/getproduct/pkg/product/provider/barcodeList"
)

func TestProductRepository(t *testing.T) {

	pr := New()
	_, err := pr.Get(context.Background(), "fake")
	if err == nil {
		t.Errorf("error should be \"product providers is empty\"")
	}

	_, err = pr.Get(context.Background(), "fake")
	if err == nil {
		t.Errorf("error should be \"context deadline exceeded\"")
	}

	prod, err := pr.Get(context.Background(), "4612732330056")
	if err != nil {
		t.Errorf(err.Error())
	}

	if prod == nil {
		t.Errorf("prod should be not nil")
	}

	if prod.Name() != "КАРСУЛЕН Раствор для инъекций (100 мл)" {
		t.Errorf("name should %s, have %s", "КАРСУЛЕН Раствор для инъекций (100 мл)", prod.Name())
	}

}
