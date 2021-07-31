package vekaptek

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/korableg/getproduct/pkg/httpUtils"
	"github.com/korableg/getproduct/pkg/product"
	"github.com/korableg/getproduct/pkg/product/provider"
)

//

type Vekaptek struct{}

func init() {
	provider.Register("vekaptek", &Vekaptek{})
}

func (v *Vekaptek) GetProduct(ctx context.Context, barcode string) (p *product.Product, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("vekaptek.ru: fetching aborted, function GetProduct was paniced")
		}
	}()

	url, err := v.getUrl(ctx, barcode)
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

	if !v.verifyBarcode(doc, barcode) {
		return nil, fmt.Errorf("vekaptek.ru: product by barcode %s not found", barcode)
	}

	p = product.New(barcode, url)
	p.SetName(v.getName(doc))
	p.SetArticle(v.getArticle(doc))
	p.SetDescription(v.getDescription(doc))
	p.SetManufacturer(v.getManufacturer(doc))
	p.SetWeight(v.getWeight(doc))

	return p, nil

}

func (v *Vekaptek) getName(doc *goquery.Document) (name string) {

	doc.Find("div.col-sm-4 h1").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		name = s.Text()
		return false
	})

	return name

}

func (v *Vekaptek) getArticle(doc *goquery.Document) (article string) {

	const splitter = "Код товара:"

	doc.Find("div.col-sm-4 ul.list-unstyled li").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {

		articleRaw := s.Text()
		if strings.HasPrefix(articleRaw, splitter) {
			splittedString := strings.Split(articleRaw, splitter)
			if len(splittedString) > 1 {
				article = strings.TrimSpace(splittedString[1])
			}
			return false
		}

		return true
	})

	return article

}

func (v *Vekaptek) getDescription(doc *goquery.Document) (description string) {

	doc.Find("#tab-description").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {

		children := s.Children()
		children.EachWithBreak(func(parentIndex int, sInternal *goquery.Selection) bool {
			if strings.HasPrefix(sInternal.Text(), "Показания") || strings.HasPrefix(sInternal.Text(), "Описание") {
				if sInternal.Nodes != nil && len(sInternal.Nodes) > 0 {
					description = sInternal.Nodes[0].NextSibling.Data
					description = strings.TrimSpace(description)
				}

				return false
			}
			return true

		})

		return false
	})

	return description

}

func (v *Vekaptek) getManufacturer(doc *goquery.Document) (manufacturer string) {

	const splitter = "Производитель:"

	doc.Find("div.col-sm-4 ul.list-unstyled li").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {

		articleRaw := s.Text()
		if strings.HasPrefix(articleRaw, splitter) {
			splittedString := strings.Split(articleRaw, splitter)
			if len(splittedString) > 1 {
				manufacturer = strings.TrimSpace(splittedString[1])
			}
			return false
		}

		return true
	})

	return manufacturer

}

func (v *Vekaptek) getWeight(doc *goquery.Document) (weight float64) {

	doc.Find("#tab-description").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {

		children := s.Children()
		children.EachWithBreak(func(parentIndex int, sInternal *goquery.Selection) bool {
			if strings.HasPrefix(sInternal.Text(), "Штрих-код") {
				if sInternal.Nodes != nil && len(sInternal.Nodes) > 0 {

					reg, _ := regexp.Compile(`Вес:\s[\d.]*`)

					weightRaw := sInternal.Nodes[0].NextSibling.Data
					weightRaw = reg.FindString(weightRaw)
					weightRaw = strings.ReplaceAll(weightRaw, "Вес:", "")
					weightRaw = strings.TrimSpace(weightRaw)

					var err error

					weight, err = strconv.ParseFloat(weightRaw, 64)
					if err != nil {
						log.Println(err)
					}
				}

				return false
			}
			return true

		})

		return false
	})

	return weight

}

func (v *Vekaptek) getUrl(ctx context.Context, barcode string) (string, error) {

	const baseUrl = "https://vekaptek.ru/search/?search=%s"

	response, err := httpUtils.Get(ctx, fmt.Sprintf(baseUrl, barcode))
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
	}

	var url string

	doc.Find("div.product-layout div.product-thumb div.caption h4 a").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {

		href, _ := s.Attr("href")
		if strings.HasPrefix(href, "https://vekaptek") {
			url = href
			return false
		}

		return true
	})

	if url == "" {
		return "", fmt.Errorf("vekaptek.ru: product with barcode %s not found", barcode)
	}
	return url, nil

}

func (v *Vekaptek) verifyBarcode(doc *goquery.Document, barcode string) bool {

	var barcodeFromPage string

	doc.Find("#tab-description").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {

		children := s.Children()
		children.EachWithBreak(func(parentIndex int, sInternal *goquery.Selection) bool {
			if strings.HasPrefix(sInternal.Text(), "Штрих-код") {
				if sInternal.Nodes != nil && len(sInternal.Nodes) > 0 {

					reg, _ := regexp.Compile(`Штрих-код:\s[\d.]*`)

					barcodeFromPage = sInternal.Nodes[0].NextSibling.Data
					barcodeFromPage = reg.FindString(barcodeFromPage)
					barcodeFromPage = strings.ReplaceAll(barcodeFromPage, "Штрих-код:", "")
					barcodeFromPage = strings.TrimSpace(barcodeFromPage)

				}

				return false
			}
			return true

		})

		return false
	})

	return barcodeFromPage == barcode

}
