package nationalCatalog

import "testing"

func TestNationalCatalog(t *testing.T) {

	const barcode = "4607004892851"

	nn := &NationalCatalog{}
	prod, err := nn.GetProduct(barcode)
	if err != nil {
		t.Error(err)
	}
	_ = prod

}
