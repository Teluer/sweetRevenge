package web

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
	"net/http"
)

type TorSession struct {
	client *http.Client
}

const proxyAddr = "localhost:1080"

var currentSession *TorSession

func openNewSession(proxyAddr string) *TorSession {
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
	currentSession = &ts
	return &ts
}

func (ts TorSession) anonymousRequest(req *http.Request) *http.Response {
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

	return resp
}

func (ts TorSession) getAnonymously(target string) *http.Response {
	if ts.client == nil {
		log.Error("Tor client not initialized")
		panic("need to init client first!")
	}

	log.Info("Sending anonymous GET to url:" + target)
	resp, err := ts.client.Get(target)
	if err != nil {
		log.WithError(err).Error("Anonymous get request failed")
		panic(err)
	}
	log.Info("Received response for anonymous GET to url:" + target)

	return resp
}
