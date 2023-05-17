package target

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sweetRevenge/websites/web"
)

const testUrl = "https://ed3e-46-55-27-40.ngrok-free.app"
const customer = "Ana Popescu"
const phone = "0000000"

func TestCookies() {
	url := "https://gudvin.md/products/skladnoj-organajzer-dlya-bagazhnika-avtomobilya-car-organizer-s-3-otdeleniyami"
	resp, err := http.Get(url)

	if err != nil {
		for _, c := range resp.Cookies() {
			fmt.Println(c)
		}
	}
}

func SendTestRequest() {
	web.OpenSession("localhost:1080").CallTor(testUrl)
}

func SendTestOrder() {
	OrderItemWithCustomerAndTargetAndItemAndLink(testUrl, customer, phone, "126", "https://gudvin.md/products/naduvnoj-divan---lamzak")
}

func Server() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:80", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	buf := new(strings.Builder)
	io.Copy(buf, r.Body)
	// check errors

	fmt.Println(buf.String())
	fmt.Println(r.RemoteAddr)
	fmt.Println(r.Header.Get("Referer"))
	for _, c := range r.Cookies() {
		fmt.Println(c)
	}
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}
}
