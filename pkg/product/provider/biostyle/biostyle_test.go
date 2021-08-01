package biostyle

import (
	"context"
	"fmt"
	"testing"
)

func TestBioStyle(t *testing.T) {

	const barcode_karsulen = "4742496003740"
	const karsulen_name = "Лимоксин-25 аэрозоль флак. 200 мл."
	// const barcode_ksila = "4742496000381"
	// const ksila_name = "Ксила флак. 50 мл."
	const barcode_fake = "fake"

	ctx := context.WithValue(context.Background(), "chromedpwsaddress", "ws://localhost:3000")

	bl := &BioStyle{}
	pr, err := bl.GetProduct(ctx, barcode_karsulen)
	if err != nil {
		t.Error(err)
	}

	if pr == nil {
		t.Errorf("func GetProduct(%s) returned nil", barcode_karsulen)
	}

	if pr != nil && pr.Name() != karsulen_name {
		t.Errorf("name should %s, have %s", karsulen_name, pr.Name())
	}

	errorTextShould := "biostyle.biz: product with barcode fake not found by google"
	pr, err = bl.GetProduct(ctx, barcode_fake)
	if err.Error() != errorTextShould {
		t.Fatal(fmt.Errorf("the error should be \"%s\"", errorTextShould))
	}

}
