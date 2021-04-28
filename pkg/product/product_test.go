package product

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestProduct(t *testing.T) {

	const barcode = "442353252342342"
	const name = "test"
	const description = "TestDescription"
	const manufacturer = "TestManufacturer"
	const unit = "TestUnit"

	p := NewProduct(barcode, name, unit, description, manufacturer)

	b, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}

	p = &Product{}
	err = json.Unmarshal(b, p)
	if err != nil {
		t.Error(err)
	}

	if p.Barcode() != barcode {
		t.Error(fmt.Errorf("barcode doesn't equal %s", barcode))
	}

	if p.Name() != name {
		t.Error(fmt.Errorf("name doesn't equal %s", name))
	}

	if p.Manufacturer() != manufacturer {
		t.Error(fmt.Errorf("manufacturer doesn't equal %s", manufacturer))
	}

	if p.Description() != description {
		t.Error(fmt.Errorf("description doesn't equal %s", description))
	}

	if p.Unit() != unit {
		t.Error(fmt.Errorf("unit doesn't equal %s", unit))
	}

}
