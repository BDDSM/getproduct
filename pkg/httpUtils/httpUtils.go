package httpUtils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"net/http"
	"net/http/cookiejar"
	url2 "net/url"
	"strings"
)

func Get(ctx context.Context, url string) (*http.Response, error) {

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0")

	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}

	response, err := client.Do(request)
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

func GetUrlByYandex(ctx context.Context, barcode, site, chromeDPWSAddress string) (string, error) {

	params := url2.Values{}
	params.Add("text", fmt.Sprintf("%s site:", barcode))

	builder := strings.Builder{}
	builder.WriteString("https://www.yandex.ru/search/?")
	builder.WriteString(params.Encode())
	builder.WriteString(site)

	body, err := GetByChromedp(ctx, chromeDPWSAddress, builder.String())
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}

	var url string

	doc.Find("a.Link").EachWithBreak(func(parentIndex int, s *goquery.Selection) bool {

		href, _ := s.Attr("href")
		if !strings.HasPrefix(href, "https://www.yandex.ru") {
			url = href
		}
		return false

	})

	if url == "" {
		return "", fmt.Errorf("%s: product with barcode %s not found by yandex", site, barcode)
	}

	return url, nil

}

func GetByChromedp(ctx context.Context, chromeDPWSAddress, url string) ([]byte, error) {

	allocatorCtx, cancel := chromedp.NewRemoteAllocator(ctx, chromeDPWSAddress)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	if err := chromedp.Run(chromeCtx, chromedp.Navigate(url)); err != nil {
		return nil, err
	}

	var body string
	if err := chromedp.Run(chromeCtx, chromedp.OuterHTML("html", &body)); err != nil {
		return nil, err
	}

	return []byte(body), nil

}
