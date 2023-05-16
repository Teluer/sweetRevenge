package fetcher

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

const baseUrl = "https://999.md"

// TODO: fetch other sections too
const ladiesUrl = "https://999.md/ru/list/dating-and-greetings/i-need-a-man"

func GetLadiesPhones() []string {
	var urls, phones []string

	currentUrl := ladiesUrl
	pageNumber := 1
	for {
		page := fetch(currentUrl) // start a goroutine
		ladies, hasNextPage := parseLadiesList(page)
		if len(ladies) > 0 {
			urls = append(urls, ladies...)
		}
		//TODO: remove test condition, make 1 second pause
		if !hasNextPage || true {
			break
		}
		pageNumber++
		currentUrl = ladiesUrl + "?page=" + strconv.Itoa(pageNumber)
	}

	//todo: use goroutines, make 1 request per second
	for _, url := range urls {
		url = baseUrl + url
		ad := fetch(url)

		phones = append(phones, getPhone(ad))
	}

	return phones
}

func parseLadiesList(htmlPage *goquery.Document) ([]string, bool) {
	hasNextString := htmlPage.Find("nav.paginator").Find(".current").Next().Length() > 0

	urls := htmlPage.Find("ul.ads-list-photo").First().Children()
	items := urls.Filter("li.ads-list-photo-item").Filter(":not(.is-adsense):not(.js-booster-inline)")
	title := items.Find(".ads-list-photo-item-title").Children()
	hrefs := title.Map(func(i int, a *goquery.Selection) string {
		href, _ := a.Attr("href")
		fmt.Println(href)
		return href
	})
	return hrefs, hasNextString
}

func getPhone(htmlPage *goquery.Document) string {
	//TODO: phone MAY BE HIDDEN, handle this case
	phoneNode := htmlPage.Find("dl.js-phone-number")
	phone, _ := phoneNode.Find("a").Attr("href")
	phone = strings.TrimPrefix(phone, "tel:+373")
	fmt.Println(phone)
	return phone
}
