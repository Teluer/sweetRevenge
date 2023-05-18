package web

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
	"net/http"
)

type TorSession struct {
	client *http.Client
}

func OpenSession(proxyAddr string) *TorSession {
	log.Info("Opening TOR session")
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		log.Error("Failed to connect to TOR!")
		log.WithError(err).Fatal("Failed to connect to TOR!")
	}

	var ts TorSession
	ts.client = &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}
	log.Info("Established TOR session successfully")
	return &ts
}

func (ts TorSession) postAnonymously(req *http.Request) {
	if ts.client == nil {
		log.Error("Tor client not initialized")
		panic("need to init client first!")
	}
	log.Info("Sending anonymous POST request to url:" + req.URL.String())
	resp, err := ts.client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("failed to send post request %v", err))
	}
	log.Info("Anonymous POST returned status " + resp.Status)
}

func (ts TorSession) GetAnonymously(target string) (*goquery.Document, []*http.Cookie) {
	if ts.client == nil {
		log.Error("Tor client not initialized")
		panic("need to init client first!")
	}

	log.Info("Sending anonymous GET to url:" + target)
	r, err := ts.client.Get(target)
	if err != nil {
		log.WithError(err).Error("Anonymous get request failed")
		panic(err)
	}
	log.Info("Received response for anonymous GET to url:" + target)

	doc, err := goquery.NewDocumentFromReader(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.WithError(err).Error("Failed to parse response")
		panic(err)
	}

	return doc, r.Cookies()
}
