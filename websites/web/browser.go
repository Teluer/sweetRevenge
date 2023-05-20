package web

import (
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
	log.Info("Got status", resp.Status, " to unsafe GET to url:", url)

	return extractDocumentFromResponse(resp)
}

func SendRequest(req *http.Request, sameSession bool) (response *http.Response, body []byte) {
	if sameSession {
		response = currentSession.anonymousRequest(req)
	} else {
		response = openNewSession().anonymousRequest(req)
	}
	return response, extractResponseBody(response)
}

func GetUrl(url string, sameSession bool) *goquery.Document {
	if sameSession {
		return extractDocumentFromResponse(currentSession.getAnonymously(url))
	} else {
		return extractDocumentFromResponse(openNewSession().getAnonymously(url))
	}
}

func GetRequest(req *http.Request, sameSession bool) *goquery.Document {
	resp, _ := SendRequest(req, sameSession)
	return extractDocumentFromResponse(resp)
}

func FetchCookies(url string, sameSession bool) []*http.Cookie {
	var resp *http.Response
	if sameSession {
		resp = currentSession.getAnonymously(url)
	} else {
		resp = openNewSession().getAnonymously(url)
	}
	resp.Body.Close()
	return resp.Cookies()
}

func extractDocumentFromResponse(resp *http.Response) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.WithError(err).Error("Failed to parse response body")
		panic(fmt.Sprintf("while parsing response body: %v", err))
	}
	resp.Body.Close() // don't leak resources
	return doc
}

func extractResponseBody(resp *http.Response) []byte {
	//TODO: handle error
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return body
}
