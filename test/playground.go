package test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sweetRevenge/websites/target"
	"sync"
)

const testUrl = "https://ed3e-46-55-27-40.ngrok-free.app"
const customer = "Ana Popescu"
const phone = "0000000"

func SendTestOrder() {
	target.OrderItemWithCustomerAndTargetAndItemAndLink(testUrl,
		customer,
		phone,
		"126",
		"https://gudvin.md/products/naduvnoj-divan---lamzak")
}

func Server() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:80", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	buf := new(strings.Builder)
	io.Copy(buf, r.Body)

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

func TestAnonSending() {
	var wg sync.WaitGroup
	wg.Add(1)
	go Server()
	SendTestOrder()
	wg.Wait()
}

//func SendTestRequest() {
//	web.openSession("localhost:1080").getAnonymously(testUrl)
//}
//
//func TestRandomCustomers() {
//	for i := 0; i < 20; i++ {
//		name, phone := target.createRandomCustomer()
//		fmt.Println(name, "  ", phone)
//	}
//
//}
