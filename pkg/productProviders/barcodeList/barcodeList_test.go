package barcodeList

import (
	"fmt"
	"github.com/korableg/getproduct/pkg/productProvider"
	"testing"
)

func TestBarcodeList(t *testing.T) {

	const barcode_karsulen = "4612732330056"
	const karsulen_name = "КАРСУЛЕН Раствор для инъекций (100 мл)"
	const barcode_ksila = "4742496000381"
	const ksila_name = "КСИЛА Раствор для инъекций (50 мл) Interchemie"
	const barcode_fake = "fake"

	bl := &BarcodeList{}
	pr, err := bl.GetProduct(barcode_karsulen)
	if err != nil {
		t.Fatal(err)
	}

	if pr == nil {
		t.Errorf("func GetProduct(%s) returned nil", barcode_karsulen)
	}

	if pr.Name() != karsulen_name {
		t.Errorf("name should %s, have %s", karsulen_name, pr.Name())
	}

	pr, err = bl.GetProduct(barcode_ksila)
	if err != nil {
		t.Fatal(err)
	}

	if pr == nil {
		t.Errorf("func GetProduct(%s) returned nil", barcode_ksila)
	}

	if pr.Name() != ksila_name {
		t.Errorf("name should %s, have %s", ksila_name, pr.Name())
	}

	pr, err = bl.GetProduct(barcode_fake)
	if err != productProvider.ErrProductDidntFind {
		t.Fatal(fmt.Errorf("the error should be \"productProvider.ErrProductDidntFind\""))
	}

}
