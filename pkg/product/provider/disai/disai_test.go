package disai

import (
	"context"
	"fmt"
	"testing"
)

func TestDisai(t *testing.T) {

	const barcode_perkutan = "4603720810186"
	const perkutan_name = "Перкутан"
	const barcode_fake = "fake"

	ctx := context.Background()

	bl := &Disai{}

	pr, err := bl.GetProduct(ctx, barcode_perkutan)
	if err != nil {
		t.Fatal(err)
	}

	if pr == nil {
		t.Errorf("func GetProduct(%s) returned nil", barcode_perkutan)
	}

	if pr.Name() != perkutan_name {
		t.Errorf("name should %s, have %s", perkutan_name, pr.Name())
	}

	errorTextShould := "biostyle.biz: product with barcode fake not found by google"
	pr, err = bl.GetProduct(ctx, barcode_fake)
	if err.Error() != errorTextShould {
		t.Fatal(fmt.Errorf("the error should be \"%s\"", errorTextShould))
	}

}
