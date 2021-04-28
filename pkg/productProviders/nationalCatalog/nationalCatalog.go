package nationalCatalog

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/korableg/getproduct/pkg/product"
	"golang.org/x/net/html"
)

const baseUrl = "https://национальный-каталог.рф"
const searchTemplate = baseUrl + "/search/?q=%s&type=goods"

type NationalCatalog struct{}

func (n *NationalCatalog) GetProduct(barcode string) (*product.Product, error) {

	resp, err := n.httpRequest(fmt.Sprintf(searchTemplate, barcode))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return n.getProduct(resp.Body)

}

func (n *NationalCatalog) httpRequest(url string) (*http.Response, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request completed with status %d", resp.StatusCode)
	}

	return resp, nil

}

func (n *NationalCatalog) getProduct(body io.Reader) (*product.Product, error) {
	htmlTokens := html.NewTokenizer(body)
	for {
		token := htmlTokens.Next()

		switch token {
		case html.ErrorToken:
			return nil, errors.New("product not found")
		case html.StartTagToken:
			{
				t := htmlTokens.Token()
				if t.Data == "a" {
					for _, v := range t.Attr {
						if v.Key == "href" && strings.HasPrefix(strings.ToLower(v.Val), "/product") {
							resp, err := n.httpRequest(baseUrl + v.Val)
							if err != nil {
								return nil, err
							}
							defer resp.Body.Close()
							return n.getProductDetails(resp.Body)
						}
					}
				}
			}

		}

	}
}

func (n *NationalCatalog) getProductDetails(body io.Reader) (*product.Product, error) {

	var barcode string
	var name string
	var unit string
	var description string
	var manufacturer string

	var nextTextIsName bool = false

	htmlTokens := html.NewTokenizer(body)
	for {
		token := htmlTokens.Next()

		switch token {
		case html.ErrorToken:
			return product.NewProduct(barcode, name, unit, description, manufacturer), nil
		case html.TextToken:
			{
				text := strings.Trim(string(htmlTokens.Text()), " ")
				lowerText := strings.ToLower(text)

				if nextTextIsName {
					name = text
					continue
				}

				if strings.HasPrefix(lowerText, "полное наименование товара") {
					nextTextIsName = true
				}

			}

		}

	}

}
