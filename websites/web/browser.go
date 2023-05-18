package web

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func SendRequest(req *http.Request, sameSession bool) *http.Response {
	if sameSession {
		return currentSession.anonymousRequest(req)
	} else {
		return openNewSession().anonymousRequest(req)
	}
}

func GetUrl(url string, sameSession bool) *goquery.Document {
	if sameSession {
		return extractDocumentFromResponse(currentSession.getAnonymously(url))
	} else {
		return extractDocumentFromResponse(openNewSession().getAnonymously(url))
	}
}

func GetRequest(req *http.Request, sameSession bool) *goquery.Document {
	return extractDocumentFromResponse(SendRequest(req, sameSession))
}

func FetchCookies(url string, sameSession bool) []*http.Cookie {
	if sameSession {
		return currentSession.getAnonymously(url).Cookies()
	} else {
		return openNewSession().getAnonymously(url).Cookies()
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
