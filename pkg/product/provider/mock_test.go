package provider

import (
	"context"
	"testing"
)

func TestProvider(t *testing.T) {

	const barcode = "111"
	const barcodeError = "222"

	var pp ProductProvider = &mockProductProvider{}

	prod, err := pp.GetProduct(context.Background(), barcode)
	if prod.Barcode() != barcode {
		t.Errorf("barcode doesn't equal %s", barcode)
	}

	if err != nil {
		t.Error(err)
	}

	_, err = pp.GetProduct(context.Background(), barcodeError)
	if err == nil {
		t.Errorf("error with barcode %s doesn't equal nil", barcodeError)
	}

}
