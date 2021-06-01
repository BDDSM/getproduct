package vekaptek

import (
	"context"
	"fmt"
	"testing"
)

func TestVekAptek(t *testing.T) {

	const barcode_karsulen = "4007221014249"
	const karsulen_name = "Катозал раствор для инъекций 10%, флакон 100 мл (вет)"
	const barcode_fake = "1234567890123"

	ctx := context.Background()

	bl := &Vekaptek{}
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

	errorTextShould := "vekaptek.ru: product with barcode 1234567890123 not found"
	pr, err = bl.GetProduct(ctx, barcode_fake)
	if err.Error() != errorTextShould {
		t.Fatal(fmt.Errorf("the error should be \"%s\"", errorTextShould))
	}

}
