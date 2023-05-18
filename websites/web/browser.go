package web

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Post(req *http.Request, sameSession bool) *http.Response {
	if sameSession {
		return currentSession.anonymousRequest(req)
	} else {
		return openNewSession(proxyAddr).anonymousRequest(req)
	}
}

func GetUrl(url string) *goquery.Document {
	return extractDocumentFromResponse(openNewSession(proxyAddr).getAnonymously(url))
}

func GetRequest(req *http.Request) *goquery.Document {
	return extractDocumentFromResponse(openNewSession(proxyAddr).anonymousRequest(req))
}

func FetchCookies(url string) []*http.Cookie {
	return openNewSession(proxyAddr).getAnonymously(url).Cookies()
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
