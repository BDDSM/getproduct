package nationalCatalog

import (
	"context"
	"fmt"
	"testing"
)

func TestNationalCatalog(t *testing.T) {

	const barcode = "4607004892851"
	const name = "Сыр плавленый Hochland Бистро Чеддер, 344 г"
	const barcode_fake = "fake"

	ctx := context.Background()

	bl := &NationalCatalog{}

	pr, err := bl.GetProduct(ctx, barcode)
	if err != nil {
		t.Fatal(err)
	}

	if pr == nil {
		t.Errorf("func GetProduct(%s) returned nil", barcode)
	}

	if pr.Name() != name {
		t.Errorf("name should %s, have %s", name, pr.Name())
	}

	errorTextShould := "национальный-каталог.рф: product with barcode fake not found by yandex"
	pr, err = bl.GetProduct(ctx, barcode_fake)
	if err.Error() != errorTextShould {
		t.Fatal(fmt.Errorf("the error should be \"%s\"", errorTextShould))
	}

}
