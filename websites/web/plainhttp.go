package web

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func sendRequest(req *http.Request) *http.Response {
	log.Info("Sending plain POST request to url:" + req.URL.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to send plain POST request")
		panic(fmt.Sprintf("failed to send post request %v", err))
	}
	log.Info("Plain post returned status " + resp.Status)
	return resp
}

func get(url string) *http.Response {
	log.Info("Sending plain GET to url:" + url)

	resp, err := http.Get(url)
	if err != nil {
		log.WithError(err).Error("Failed to send plain GET to url:" + url)
		panic(fmt.Sprintf("while reading %s: %v", url, err))
	}
	return resp
}
