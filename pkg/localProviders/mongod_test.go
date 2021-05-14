package localProviders

import (
	"context"
	"testing"
)

const hostname = "localhost"
const port = 27017
const username = ""
const password = ""

const barcode = "4612732330056"

func TestNew(t *testing.T) {

	m, err := New(hostname, port, username, password)
	if err != nil {
		t.Error(err)
	}
	_ = m

}

func TestGetProduct(t *testing.T) {

	ctx := context.Background()

	m, err := New(hostname, port, username, password)
	if err != nil {
		t.Error(err)
	}
	p, err := m.GetProduct(ctx, barcode)
	if err != nil {
		t.Error(err)
	}

	_ = p

}
