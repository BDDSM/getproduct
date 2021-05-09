package barcodeList

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/korableg/getproduct/pkg/httpUtils"
	"github.com/korableg/getproduct/pkg/product"
	"strconv"
	"strings"
)

const endpointTemplate = "https://barcode-list.ru/barcode/RU/Поиск.htm?barcode=%s"

type BarcodeList struct{}

func (b *BarcodeList) GetProduct(ctx context.Context, barcode string) (*product.Product, error) {

	url := fmt.Sprintf(endpointTemplate, barcode)

	response, err := httpUtils.Get(ctx, url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	name := getName(doc)
	if name == "" {
		return nil, fmt.Errorf("barcode-list.ru: product by barcode %s not found", barcode)
	}
	unit := getUnit(doc)

	p := product.New(barcode, url)
	p.SetName(name)
	p.SetUnit(unit)

	return p, nil

}

func getName(doc *goquery.Document) (name string) {

	doc.Find(".pageTitle").EachWithBreak(func(i int, s *goquery.Selection) bool {
		//if product didn't find .pageTitle has a "поиск:" text.
		//In this case func returns false without filled properties
		if strings.HasPrefix(strings.ToLower(s.Text()), "поиск") {
			return false
		}
		name = prepareName(s.Text())
		return false
	})

	return name

}

func prepareName(title string) string {

	splittedTitle := strings.Split(title, " - Штрих-код")
	if len(splittedTitle) > 0 {
		return splittedTitle[0]
	}

	return title

}

func getUnit(doc *goquery.Document) (unit string) {

	var maxRate = 0
	var tempUnit = ""
	doc.Find(".randomBarcodes tr").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		if parentIndex == 0 {
			return true
		}

		s.Children().Each(func(childIndex int, td *goquery.Selection) {
			switch childIndex {
			case 3:
				tempUnit = td.Text()
			case 4:
				{
					if rate, err := strconv.Atoi(td.Text()); err == nil && rate > maxRate {
						unit = tempUnit
						maxRate = rate
					}
				}

			}
		})

		return true
	})

	return unit

}
