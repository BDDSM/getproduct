package httpUtils

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	url2 "net/url"
	"strings"
)

func Get(ctx context.Context, url string) (*http.Response, error) {

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request with url: \"%s\" finished with status: %d", url, response.StatusCode)
	}

	return response, nil

}

func GetUrlByGoogle(ctx context.Context, barcode, site string) (string, error) {

	params := url2.Values{}
	params.Add("q", fmt.Sprintf("%s site:%s", barcode, site))

	builder := strings.Builder{}
	builder.WriteString("https://www.google.com/search?")
	builder.WriteString(params.Encode())

	response, err := Get(ctx, builder.String())
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
	}

	var url string

	doc.Find("div.g div div div a").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {

		href, _ := s.Attr("href")
		url = href
		return false

	})

	if url == "" {
		return "", fmt.Errorf("%s: product with barcode %s not found by google", site, barcode)
	}

	return url, nil

}
