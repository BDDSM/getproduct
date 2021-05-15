package localProviders

import (
	"context"
	"testing"

	"github.com/korableg/getproduct/pkg/product"
)

const hostname = "localhost"
const port = 27017
const username = ""
const password = ""

const barcode = "testBarcode"

const collectionNameTest = "products_test"

func TestMongoDB(t *testing.T) {

	ctx := context.Background()

	m, err := newMongo(hostname, port, username, password, collectionNameTest)
	if err != nil {
		t.Error(err)
	}

	p, err := m.getProduct(ctx, barcode, collectionNameTest)
	if err != nil {
		t.Error(err)
	}

	if p != nil {
		t.Error("product should be nil")
	}

	p = product.New(barcode, "testURL")
	p.SetName("Test")
	p.SetArticle("235435")
	p.SetUnit("шт")

	err = m.deleteProduct(ctx, p, collectionNameTest)
	if err != nil {
		t.Error(err)
	}

	err = m.addProduct(ctx, p, collectionNameTest)
	if err != nil {
		t.Error(err)
	}

	p, err = m.getProduct(ctx, barcode, collectionNameTest)
	if err != nil {
		t.Error(err)
	}
	if p == nil {
		t.Error("product should be not nil")
	} else if p.Barcode() != barcode {
		t.Errorf("barcode should be %s, have %s", barcode, p.Barcode())
	}

	p = product.New(barcode, "testURL")
	p.SetName("Test1")
	p.SetArticle("235435")
	p.SetUnit("шт")

	err = m.addProduct(ctx, p, collectionNameTest)
	if err != nil {
		t.Error(err)
	}

	p, err = m.getProduct(ctx, barcode, collectionNameTest)
	if err != nil {
		t.Error(err)
	}

	if p == nil {
		t.Error("product should be not nil")
	} else if p.Name() != "Test1" {
		t.Errorf("product hasn't replaced with barcode %s", barcode)
	}

	err = m.deleteProduct(ctx, p, collectionNameTest)
	if err != nil {
		t.Error(err)
	}

}
