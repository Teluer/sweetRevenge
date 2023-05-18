package web

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

const proxyAddr = "localhost:1080"

func Post(req *http.Request, anon bool) {
	if anon {
		OpenSession(proxyAddr).postAnonymously(req)
	} else {
		post(req)
	}
}

func Fetch(url string, anon bool) *goquery.Document {
	result, _ := FetchWithCookies(url, anon)
	return result
}

func FetchWithCookies(url string, anon bool) (*goquery.Document, []*http.Cookie) {
	if anon {
		return OpenSession(proxyAddr).GetAnonymously(url)
	} else {
		return get(url)
	}
}
