package fetcher

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

// TODO: use browser method instead
func fetch(url string) *goquery.Document {
	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("while reading %s: %v", url, err))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("while parsing %s: %v", url, err))
	}
	resp.Body.Close() // don't leak resources

	return doc
}
