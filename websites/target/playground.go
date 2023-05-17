package target

import (
	"fmt"
	"net/http"
)

func TestCookies() {
	url := "https://gudvin.md/products/skladnoj-organajzer-dlya-bagazhnika-avtomobilya-car-organizer-s-3-otdeleniyami"
	resp, err := http.Get(url)

	if err != nil {
		for _, c := range resp.Cookies() {
			fmt.Println(c)
		}
	}
}
