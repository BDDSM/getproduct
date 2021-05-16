package biostyle

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/korableg/getproduct/pkg/httpUtils"
	"github.com/korableg/getproduct/pkg/product"
)

type BioStyle struct {
	chromeDPWSAddress string
}

func New(chromeDPWSAddress string) *BioStyle {
	b := BioStyle{
		chromeDPWSAddress: chromeDPWSAddress,
	}

	return &b

}

func (b *BioStyle) GetProduct(ctx context.Context, barcode string) (p *product.Product, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("biostyle.biz: fetching aborted, function GetProduct was paniced")
		}
	}()

	url, err := httpUtils.GetUrlByGoogle(ctx, barcode, "biostyle.biz")
	if err != nil {
		return nil, err
	}

	body, err := httpUtils.GetByChromedp(ctx, b.chromeDPWSAddress, url)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	if !b.verifyBarcode(doc, barcode) {
		return nil, fmt.Errorf("biostyle.biz: product by barcode %s not found", barcode)
	}

	p = product.New(barcode, url)
	p.SetName(b.getName(doc))
	p.SetArticle(b.getArticle(doc))
	p.SetDescription(b.getDescription(doc))
	p.SetManufacturer(b.getManufacturer(doc))
	p.SetUnit(b.getUnit(doc))
	p.SetWeight(b.getWeight(doc))
	p.SetPicture(b.getPicture(ctx, doc))

	return p, nil

}

func (b *BioStyle) getName(doc *goquery.Document) (name string) {
	doc.Find("#pagetitle").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		name = s.Text()
		return false
	})
	return name
}

func (b *BioStyle) getArticle(doc *goquery.Document) (article string) {
	doc.Find("span.article__value").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		article = s.Text()
		return false
	})
	return article
}

func (b *BioStyle) getDescription(doc *goquery.Document) (description string) {
	doc.Find("div.content[itemprop=\"description\"]").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		description = s.Text()
		return false
	})

	description = strip.StripTags(description)
	description = strings.TrimSpace(description)

	return description
}

func (b *BioStyle) getManufacturer(doc *goquery.Document) (manufacturer string) {
	doc.Find("a.brand__link").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		manufacturer = s.Text()
		return false
	})

	return manufacturer
}

func (b *BioStyle) getUnit(doc *goquery.Document) (unit string) {

	unit = b.getAdditionalProperty(doc, "базовая единица")

	return unit
}

func (b *BioStyle) getWeight(doc *goquery.Document) (weight float64) {

	weightRaw := b.getAdditionalProperty(doc, "вес")
	weight, err := strconv.ParseFloat(weightRaw, 64)
	if err != nil {
		log.Println(err)
	}

	return weight
}

func (b *BioStyle) getPicture(ctx context.Context, doc *goquery.Document) (picture []byte) {

	doc.Find("img.product-detail-gallery__picture").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		if src, ok := s.Attr("src"); ok {
			if !strings.HasPrefix(src, "https") {
				src = "https:" + src
			}
			response, err := httpUtils.Get(ctx, src)
			if err != nil {
				log.Println(err)
			}
			defer response.Body.Close()
			picture, err = io.ReadAll(response.Body)
		}

		return false
	})
	return picture

}

func (b *BioStyle) verifyBarcode(doc *goquery.Document, barcode string) bool {
	barcodeFromPage := b.getAdditionalProperty(doc, "штрихкод")
	return barcodeFromPage == barcode
}

func (b *BioStyle) getAdditionalProperty(doc *goquery.Document, key string) (value string) {

	doc.Find("tr[itemprop=\"additionalProperty\"]").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {

		itemFound := false
		s.Find("span[itemprop=\"name\"]").EachWithBreak(func(parentIndex int, sInternal *goquery.Selection) bool {
			if strings.HasPrefix(strings.ToLower(sInternal.Text()), key) {
				itemFound = true
			}
			return false
		})

		if itemFound {
			s.Find("span[itemprop=\"value\"]").EachWithBreak(func(parentIndex int, sInternal *goquery.Selection) bool {
				value = sInternal.Text()
				return false
			})
		}

		return !itemFound

	})

	value = strings.TrimSpace(value)

	return value

}
