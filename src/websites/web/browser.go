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

func SendRequest(req *http.Request, sameSession bool) (response *http.Response, body []byte) {
	if sameSession {
		return currentSession.anonymousRequest(req)
	} else {
		return openNewSession().anonymousRequest(req)
	}
}

func GetUrl(url string, sameSession bool) *goquery.Document {
	var responseBody []byte
	if sameSession {
		_, responseBody = currentSession.getAnonymously(url)
	} else {
		_, responseBody = openNewSession().getAnonymously(url)
	}
	return extractDocumentFromResponseBody(responseBody)
}

func GetRequest(req *http.Request, sameSession bool) *goquery.Document {
	_, body := SendRequest(req, sameSession)
	return extractDocumentFromResponseBody(body)
}

func FetchCookies(url string, sameSession bool) []*http.Cookie {
	var resp *http.Response
	if sameSession {
		resp, _ = currentSession.getAnonymously(url)
	} else {
		resp, _ = openNewSession().getAnonymously(url)
	}
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
