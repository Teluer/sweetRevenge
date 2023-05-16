package tor

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

type TorSession struct {
	client *http.Client
}

type Browser struct {
	torSession *TorSession
}

func (b Browser) Fetch(url string, anon bool) *goquery.Document {
	if anon {
		return b.fetchAnonimously(url)
	} else {
		return b.fetchUnsafe(url)
	}
}

func (b Browser) fetchUnsafe(url string) *goquery.Document {
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

//TODO: init session here, always open new session on request.
func (b Browser) fetchAnonimously(url string) *goquery.Document {
	return b.torSession.callTor(url)
}

func (ts TorSession) Init(proxyAddr string) {
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

	ts.client = &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}
}

// TODO: see if auth is needed
// TODO: return something better than string
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
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()
	return doc
}
