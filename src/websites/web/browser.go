package web

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func GetUnsafe(url string) *goquery.Document {
	resp, err := http.Get(url)
	if err != nil {
		log.WithError(err).Error("Unsafe GET request failed")
		panic(err)
	}
	defer resp.Body.Close()
	log.Info("Got status", resp.Status, " to unsafe GET to url:", url)

	body, _ := io.ReadAll(resp.Body)
	return extractDocumentFromResponseBody(body)
}

func SendRequest(req *http.Request) (response *http.Response, body []byte) {
	return openNewSession().anonymousRequest(req)
}

func GetUrl(url string) *goquery.Document {
	_, responseBody := openNewSession().getAnonymously(url)
	return extractDocumentFromResponseBody(responseBody)
}

func GetRequest(req *http.Request) *goquery.Document {
	_, body := SendRequest(req)
	return extractDocumentFromResponseBody(body)
}

func FetchCookies(url string) []*http.Cookie {
	var resp *http.Response
	resp, _ = openNewSession().getAnonymously(url)
	return resp.Cookies()
}

func extractDocumentFromResponseBody(body []byte) *goquery.Document {
	reader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.WithError(err).Error("Failed to parse response body")
		panic(fmt.Sprintf("while parsing response body: %v", err))
	}

	return doc
}
