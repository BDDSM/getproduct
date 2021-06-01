package biostyle

import (
	"context"
	"fmt"
	"testing"
)

func TestBioStyle(t *testing.T) {

	const barcode_karsulen = "4612732330056"
	const karsulen_name = "Карсулен флакон, 100 мл"
	const barcode_ksila = "4742496000381"
	const ksila_name = "Ксила флак. 50 мл."
	const barcode_fake = "fake"

	ctx := context.Background()

	bl := New("ws://localhost:3000")
	pr, err := bl.GetProduct(ctx, barcode_karsulen)
	if err != nil {
		t.Fatal(err)
	}

	if pr == nil {
		t.Errorf("func GetProduct(%s) returned nil", barcode_karsulen)
	}

	if pr.Name() != karsulen_name {
		t.Errorf("name should %s, have %s", karsulen_name, pr.Name())
	}

	pr, err = bl.GetProduct(ctx, barcode_ksila)
	if err != nil {
		t.Fatal(err)
	}

	if pr == nil {
		t.Errorf("func GetProduct(%s) returned nil", barcode_ksila)
	}

	if pr.Name() != ksila_name {
		t.Errorf("name should %s, have %s", ksila_name, pr.Name())
	}

	errorTextShould := "biostyle.biz: product with barcode fake not found by google"
	pr, err = bl.GetProduct(ctx, barcode_fake)
	if err.Error() != errorTextShould {
		t.Fatal(fmt.Errorf("the error should be \"%s\"", errorTextShould))
	}

}
