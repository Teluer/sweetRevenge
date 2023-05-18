package web

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const proxyAddr = "localhost:1080"

func Post(req *http.Request, anon bool) *http.Response {
	if anon {
		return openSession(proxyAddr).anonymousRequest(req)
	} else {
		return sendRequest(req)
	}
}

func GetUrl(url string, anon bool) (doc *goquery.Document) {
	if anon {
		return extractDocumentFromResponse(openSession(proxyAddr).getAnonymously(url))
	} else {
		return extractDocumentFromResponse(get(url))
	}
	return doc
}

func GetRequest(req *http.Request, anon bool) *goquery.Document {
	if anon {
		return extractDocumentFromResponse(openSession(proxyAddr).anonymousRequest(req))
	} else {
		return extractDocumentFromResponse(sendRequest(req))
	}
}

func FetchCookies(url string, anon bool) []*http.Cookie {
	if anon {
		return openSession(proxyAddr).getAnonymously(url).Cookies()
	} else {
		return get(url).Cookies()
	}
}

func extractDocumentFromResponse(resp *http.Response) (doc *goquery.Document) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.WithError(err).Error("Failed to parse response body")
		panic(fmt.Sprintf("while parsing response body: %v", err))
	}
	resp.Body.Close() // don't leak resources
	return doc
}
