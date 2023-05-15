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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/proxy"
)

// TODO: replace example.org with parameter, see if auth is needed
// TODO: return something better than string
func tor() {
	target := flag.String("target", "http://example.org", "URL to get")
	proxyAddr := flag.String("proxy", "localhost:1080", "SOCKS5 proxy address to use")
	//username := flag.String("user", "", "username for SOCKS5 proxy")
	//password := flag.String("pass", "", "password for SOCKS5 proxy")
	flag.Parse()

	//auth := proxy.Auth{
	//	User:     *username,
	//	Password: *password,
	//}
	//dialer, err := proxy.SOCKS5("tcp", *proxyAddr, &auth, nil)
	dialer, err := proxy.SOCKS5("tcp", *proxyAddr, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}

	r, err := client.Get(*target)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
