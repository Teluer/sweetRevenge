package web

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

func Post(req *http.Request, anon bool) {
	if anon {
		panic("NOT TESTED YET")
		openSession(proxyAddr).postAnonymously(req)
	} else {
		postUnsafe(req)
	}
}

func Fetch(url string, anon bool) *goquery.Document {
	result, _ := FetchWithCookies(url, anon)
	return result
}

func FetchWithCookies(url string, anon bool) (*goquery.Document, []*http.Cookie) {
	if anon {
		panic("NOT TESTED YET")
		return openSession(proxyAddr).callTor(url)
	} else {
		return fetchUnsafe(url)
	}
}

func fetchUnsafe(url string) (*goquery.Document, []*http.Cookie) {
	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("while reading %s: %v", url, err))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("while parsing %s: %v", url, err))
	}
	resp.Body.Close() // don't leak resources
	return doc, resp.Cookies()
}

func postUnsafe(req *http.Request) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("failed to send post request %v", err))
	}
	fmt.Printf("Post returned status %s", resp.Status)
}

func (ts TorSession) postAnonymously(req *http.Request) {
	if ts.client == nil {
		panic("need to init client first!")
	}
	resp, err := ts.client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("failed to send post request %v", err))
	}
	fmt.Printf("Post returned status %s", resp.Status)
}

func openSession(proxyAddr string) *TorSession {
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
func (ts TorSession) callTor(target string) (*goquery.Document, []*http.Cookie) {
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

	return doc, r.Cookies()
}
