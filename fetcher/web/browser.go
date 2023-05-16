package web

/*
$ go run http-get-socks-transport.go -proxy localhost:1080 \
    -user myuser -pass mypass \
    http://example.org
<!doctype html>
<html>
// ... rest of response
*/

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"

	"golang.org/x/net/proxy"
)

const proxyAddr = "localhost:1080"

type TorSession struct {
	client *http.Client
}

type Browser struct {
	torSession *TorSession
}

var b = Browser{initialize(proxyAddr)}

func Fetch(url string, anon bool) *goquery.Document {
	if anon {
		return fetchAnonimously(url)
	} else {
		return fetchUnsafe(url)
	}
}

func fetchUnsafe(url string) *goquery.Document {
	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("while reading %s: %v", url, err))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("while parsing %s: %v", url, err))
	}
	resp.Body.Close() // don't leak resources
	return doc
}

// TODO: init session here, always open new session on request.
func fetchAnonimously(url string) *goquery.Document {
	return b.torSession.callTor(url)
}

func initialize(proxyAddr string) *TorSession {
	//proxyAddr := flag.String("proxy", "localhost:1080", "SOCKS5 proxy address to use")
	//username := flag.String("user", "", "username for SOCKS5 proxy")
	//password := flag.String("pass", "", "password for SOCKS5 proxy")
	//flag.Parse()

	//auth := proxy.Auth{
	//	User:     *username,
	//	Password: *password,
	//}
	//dialer, err := proxy.SOCKS5("tcp", *proxyAddr, &auth, nil)
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, nil)
	if err != nil {
		//TODO implement logging and error handling
		log.Fatal(err)
	}

	var ts TorSession
	ts.client = &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}
	return &ts
}

// TODO: see if auth is needed
func (ts TorSession) callTor(target string) *goquery.Document {
	//target := flag.String("target", "http://example.org", "URL to get")
	if ts.client == nil {
		panic("need to init client first!")
	}

	r, err := ts.client.Get(target)
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}

	return doc
}
