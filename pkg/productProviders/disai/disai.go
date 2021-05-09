package disai

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/korableg/getproduct/pkg/httpUtils"
	"github.com/korableg/getproduct/pkg/product"
	"golang.org/x/net/html/charset"
	"strings"
)

type Disai struct{}

func (d *Disai) GetProduct(ctx context.Context, barcode string) (*product.Product, error) {

	url := fmt.Sprintf("http://ru.disai.org/barcode/ean-13/%s", barcode)

	response, err := httpUtils.Get(ctx, url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	reader, err := charset.NewReader(response.Body, "text/html; charset=windows-1251")
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	name := d.getName(doc)
	if name == "" {
		return nil, fmt.Errorf("ru.disai.org: product by barcode %s not found", barcode)
	}

	p := product.New(barcode, url)
	p.SetName(name)
	p.SetManufacturer(d.getManufacturer(doc))

	return p, nil

}

func (d *Disai) getName(doc *goquery.Document) (name string) {

	doc.Find("div.caption h1").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		name = s.Text()
		return false
	})

	name = strings.TrimSpace(name)

	return name

}

func (d *Disai) getManufacturer(doc *goquery.Document) (manufacturer string) {

	doc.Find("div.caption p a font").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		manufacturer = s.Text()
		return false
	})

	manufacturer = strings.TrimSpace(manufacturer)

	return manufacturer

}
