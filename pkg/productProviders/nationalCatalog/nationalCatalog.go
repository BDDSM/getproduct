package nationalCatalog

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/korableg/getproduct/pkg/httpUtils"
	"github.com/korableg/getproduct/pkg/product"
	"io"
	"log"
	"strconv"
	"strings"
)

type NationalCatalog struct {
	chromedpWsAddress string
}

func New(chromedpWsAddress string) *NationalCatalog {
	nc := NationalCatalog{
		chromedpWsAddress: chromedpWsAddress,
	}

	return &nc
}

func (nc *NationalCatalog) GetProduct(ctx context.Context, barcode string) (*product.Product, error) {
	url, err := httpUtils.GetUrlByYandex(ctx, barcode, "национальный-каталог.рф", nc.chromedpWsAddress)
	if err != nil {
		return nil, err
	}

	response, err := httpUtils.Get(ctx, url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	if !nc.verifyBarcode(doc, barcode) {
		return nil, fmt.Errorf("национальный-каталог.рф: product by barcode %s not found", barcode)
	}

	p := product.New(barcode, url)
	p.SetName(nc.getName(doc))
	p.SetUnit(nc.getUnit(doc))
	p.SetWeight(nc.getWeight(doc))
	p.SetDescription(nc.getDescription(doc))
	p.SetPicture(nc.getPicture(ctx, doc))

	return p, nil

}

func (nc *NationalCatalog) getName(doc *goquery.Document) (name string) {

	doc.Find("div.container h1").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		name = s.Text()
		return false
	})

	name = strings.TrimSpace(name)

	return name

}

func (nc *NationalCatalog) getUnit(doc *goquery.Document) (unit string) {

	unit = nc.getPropertyFromTable(doc, "базовая единица")

	return unit
}

func (nc *NationalCatalog) getDescription(doc *goquery.Document) (description string) {

	description = nc.getPropertyFromTable(doc, "состав")

	return description
}

func (nc *NationalCatalog) getWeight(doc *goquery.Document) (weight float64) {

	weightRaw := nc.getPropertyFromTable(doc, "заявленный объём")

	var err error

	weight, err = strconv.ParseFloat(weightRaw, 64)
	if err != nil {
		log.Println(err)
	}

	return weight

}

func (nc *NationalCatalog) getPicture(ctx context.Context, doc *goquery.Document) (picture []byte) {

	doc.Find("img.productSinglePhoto").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {

		if url, ok := s.Attr("src"); ok {
			response, err := httpUtils.Get(ctx, url)
			if err != nil {
				log.Println(err)
			}
			defer response.Body.Close()

			picture, err = io.ReadAll(response.Body)
			if err != nil {
				log.Println(err)
			}

		}

		return false
	})

	return picture

}

func (nc *NationalCatalog) getPropertyFromTable(doc *goquery.Document, key string) (value string) {

	doc.Find("table.table th").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {

		sText := strings.TrimSpace(s.Text())
		sText = strings.ToLower(sText)

		if strings.HasPrefix(sText, key) {
			if s.Nodes != nil && len(s.Nodes) > 0 {
				node := s.Nodes[0]
				for i := 0; i < 5 && node.NextSibling != nil && node.Data != "td"; i++ {
					node = node.NextSibling
				}
				if node.Data == "td" && node.FirstChild != nil {
					value = node.FirstChild.Data
				}
			}
			return false
		}
		return true
	})

	value = strings.TrimSpace(value)

	return value

}

func (nc *NationalCatalog) verifyBarcode(doc *goquery.Document, barcode string) (match bool) {

	doc.Find("div.good-datamatrix div.pull-right").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		match = s.Text() == barcode
		return false
	})

	return match
}
