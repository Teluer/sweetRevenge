package web

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func post(req *http.Request) {
	log.Info("Sending plain POST request to url:" + req.URL.String())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to send plain POST request")
		panic(fmt.Sprintf("failed to send post request %v", err))
	}
	log.Info("Plain post returned status " + resp.Status)
}

func get(url string) (*goquery.Document, []*http.Cookie) {
	log.Info("Sending plain GET to url:" + url)

	resp, err := http.Get(url)
	if err != nil {
		log.WithError(err).Error("Failed to send plain GET to url:" + url)
		panic(fmt.Sprintf("while reading %s: %v", url, err))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.WithError(err).Error("Failed to parse response to plain GET to url: " + url)
		panic(fmt.Sprintf("while parsing %s: %v", url, err))
	}
	resp.Body.Close() // don't leak resources
	return doc, resp.Cookies()
}
