package fetcher

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

const baseUrl = "https://999.md"
const ladiesUrl = "https://999.md/ru/list/dating-and-greetings/i-need-a-man"

func GetLadiesPhones() []string {
	var urls, phones []string
	//TODO: get list of ads and add to urls. if last page is current, stop.

	currentUrl := ladiesUrl
	pageNumber := 1
	for {
		page := fetch(currentUrl) // start a goroutine
		ladies, hasNextPage := parseLadiesList(page)
		if len(ladies) > 0 {
			urls = append(urls, ladies...)
		}
		//TODO
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

	//for _, url := range urls {
	//	go fetch(url, ch) // start a goroutine
	//}
	//for range urls {8
	//	fmt.Println(<-ch) // receive from channel ch
	//}
}

// TODO: get all ads - solve search bug
func parseLadiesList(htmlPage *goquery.Document) ([]string, bool) {
	//<nav class="paginator cf"> <ul>
	//<li class="current"><a href="/ru/list/dating-and-greetings/i-need-a-man" data-pjax="">
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
	//****MAY BE HIDDEN
	//<dl class="js-phone-number adPage__content__phone grid_18">
	//   <a href="tel:+37368829133">
	phoneNode := htmlPage.Find("dl.js-phone-number")
	phone, _ := phoneNode.Find("a").Attr("href")
	phone = strings.TrimPrefix(phone, "tel:+373")
	fmt.Println(phone)
	return phone
}
