package eapteka

import (
	"context"
	"fmt"
	"testing"
)

const chromedpWSAddress = "ws://localhost:3000"

func TestEapteka(t *testing.T) {

	const barcode_ksila = "8718692823822"
	const ksila_name = "Ксила флак. 50 мл."
	const barcode_fake = "fake"

	ctx := context.Background()

	bl := New(chromedpWSAddress)

	pr, err := bl.GetProduct(ctx, barcode_ksila)
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
