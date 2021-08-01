package product

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
	"sync"
)

// Product model
// swagger:model product
type Product struct {
	// Barcode of product
	// required: true
	barcode string
	// Article of product
	article string
	// Name of product
	name string
	// Description of product
	description string
	// Manufacturer of product
	manufacturer string
	// Unit of product
	unit string
	// Weight of 1 unit product
	weight float64
	// Url where was found this product
	// required: true
	url string
	// Picture of product
	picture []byte
	// Additional properties
	properties   map[string]string
	propertiesMu *sync.Mutex
}

func New(barcode string, url string) *Product {

	p := Product{
		barcode:      barcode,
		url:          url,
		properties:   make(map[string]string),
		propertiesMu: &sync.Mutex{},
	}

	return &p
}

func (p *Product) Barcode() string {
	return p.barcode
}

func (p *Product) Url() string {
	return p.url
}

func (p *Product) Article() string {
	return p.article
}

func (p *Product) SetArticle(article string) {
	p.article = article
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) SetName(name string) {
	p.name = name
}

func (p *Product) Description() string {
	return p.description
}

func (p *Product) SetDescription(description string) {
	p.description = description
}

func (p *Product) Manufacturer() string {
	return p.manufacturer
}

func (p *Product) SetManufacturer(manufacturer string) {
	p.manufacturer = manufacturer
}

func (p *Product) Unit() string {
	return p.unit
}

func (p *Product) SetUnit(unit string) {
	unit = processUnit(unit)
	p.unit = unit
}

func (p *Product) Weight() float64 {
	return p.weight
}

func (p *Product) SetWeight(weight float64) {
	p.weight = weight
}

func (p *Product) Picture() []byte {
	return p.picture
}

func (p *Product) SetPicture(picture []byte) {
	p.picture = picture
}

func (p *Product) Properties() map[string]string {
	return p.properties
}

func (p *Product) AddProperty(name, value string) {

	key := strings.TrimSpace(strings.ToLower(name))
	if len(key) == 0 {
		return
	}

	p.propertiesMu.Lock()
	defer p.propertiesMu.Unlock()

	p.properties[key] = value

}

func (p *Product) MarshalJSON() ([]byte, error) {

	prodMap := make(map[string]interface{})
	prodMap["barcode"] = p.barcode
	prodMap["article"] = p.article
	prodMap["name"] = p.name
	prodMap["description"] = p.description
	prodMap["manufacturer"] = p.manufacturer
	prodMap["unit"] = p.unit
	prodMap["weight"] = p.weight
	prodMap["url"] = p.url
	if p.picture == nil {
		prodMap["picture"] = nil
	} else {
		prodMap["picture"] = base64.StdEncoding.EncodeToString(p.picture)
	}
	prodMap["properties"] = p.properties

	return json.Marshal(prodMap)

}

func (p *Product) UnmarshalJSON(b []byte) error {

	p.properties = make(map[string]string)

	prodMap := make(map[string]interface{})
	err := json.Unmarshal(b, &prodMap)
	if err != nil {
		return err
	}

	p.barcode = prodMap["barcode"].(string)
	p.article = prodMap["article"].(string)
	p.name = prodMap["name"].(string)
	p.description = prodMap["description"].(string)
	p.manufacturer = prodMap["manufacturer"].(string)
	p.unit = prodMap["unit"].(string)
	p.weight = prodMap["weight"].(float64)
	p.url = prodMap["url"].(string)
	if prodMap["picture"] != nil {
		if picture, err := base64.StdEncoding.DecodeString(prodMap["picture"].(string)); err == nil {
			p.picture = picture
		} else {
			log.Println(err)
		}
	}
	for k, v := range prodMap["properties"].(map[string]interface{}) {
		p.properties[k] = v.(string)
	}

	return nil

}

func (p *Product) Rating() uint64 {
	var rating uint64 = 0

	rating += uint64(len(p.name))
	rating += uint64(len(p.article))
	rating += uint64(len(p.description))
	rating += uint64(len(p.manufacturer))

	if p.weight != 0.0 {
		rating++
	}

	if p.unit != "" {
		rating++
	}

	if p.picture != nil {
		rating += uint64(len(p.picture))
	}

	return rating
}

func processUnit(unit string) string {
	unit = strings.ToLower(unit)
	unit = strings.Trim(unit, ".")
	return unit
}
