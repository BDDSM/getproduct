package product

import "encoding/json"

type Product struct {
	barcode     string
	name        string
	description string
}

func (p *Product) MarshalJson() ([]byte, error) {

	prodMap := make(map[string]interface{})
	prodMap["barcode"] = p.barcode
	prodMap["name"] = p.name
	prodMap["description"] = p.description

	return json.Marshal(prodMap)

}
