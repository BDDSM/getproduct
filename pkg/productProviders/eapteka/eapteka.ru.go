package eapteka

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/korableg/getproduct/pkg/httpUtils"
	"github.com/korableg/getproduct/pkg/product"
)

type Eapteka struct{}

func (e *Eapteka) GetProduct(ctx context.Context, barcode string) (*product.Product, error) {

	url, err := httpUtils.GetUrlByGoogle(ctx, barcode, "eapteka.ru")
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

	p := product.New(barcode, url)
	p.SetName(e.getName(doc))
	p.SetArticle(e.getArticle(doc))

	return p, nil

}

func (e *Eapteka) getName(doc *goquery.Document) (name string) {
	doc.Find("div[itemtype=\"http://schema.org/Product\"] h1").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		name = s.Text()
		return false
	})
	return name
}

func (e *Eapteka) getArticle(doc *goquery.Document) (article string) {
	doc.Find("span[data-action=\"article\"]").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {
		article = s.Text()
		return false
	})
	return article
}
