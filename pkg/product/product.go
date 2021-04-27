package product

import "encoding/json"

type Product struct {
	barcode      string
	name         string
	description  string
	manufacturer string
	unit         string
}

func NewProduct(barcode, name, unit, description, manufacturer string) *Product {
	p := Product{
		barcode:      barcode,
		name:         name,
		unit:         unit,
		description:  description,
		manufacturer: manufacturer,
	}

	return &p
}

func (p *Product) Barcode() string {
	return p.barcode
}

func (p *Product) Name() string {
	return p.barcode
}

func (p *Product) Description() string {
	return p.description
}

func (p *Product) Manufacturer() string {
	return p.manufacturer
}

func (p *Product) Unit() string {
	return p.unit
}

func (p *Product) MarshalJSON() ([]byte, error) {

	prodMap := make(map[string]interface{})
	prodMap["barcode"] = p.barcode
	prodMap["name"] = p.name
	prodMap["manufacturer"] = p.manufacturer
	prodMap["description"] = p.description
	prodMap["unit"] = p.unit

	return json.Marshal(prodMap)

}

func (p *Product) UnmarshalJSON(b []byte) error {

	prodMap := make(map[string]interface{})
	err := json.Unmarshal(b, &prodMap)
	if err != nil {
		return err
	}

	p.barcode = prodMap["barcode"].(string)
	p.name = prodMap["name"].(string)
	p.unit = prodMap["unit"].(string)
	p.description = prodMap["description"].(string)
	p.manufacturer = prodMap["manufacturer"].(string)

	return nil

}
