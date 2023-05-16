package fetcher

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"sweetRevenge/db/dto"
	"sweetRevenge/web"
	"time"
)

const baseUrl = "https://999.md"

// TODO: fetch other sections too
const ladiesUrl = "https://999.md/ru/list/dating-and-greetings/i-need-a-man"

func GetLadies() (ladies []dto.Lady) {
	var urls []string
	currentUrl := ladiesUrl
	pageNumber := 1
	for {
		page := web.Fetch(currentUrl, false) // start a goroutine
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

	//TODO: use goroutines, make 1 request per second
	for _, url := range urls {
		url = baseUrl + url
		ad := web.Fetch(url, false)

		ladies = append(ladies, getLady(ad))
	}
	return ladies
}

func parseLadiesList(htmlPage *goquery.Document) (adLinks []string, hasNextPage bool) {
	hasNextPage = htmlPage.Find("nav.paginator").Find(".current").Next().Length() > 0

	adLinks = htmlPage.Find("ul.ads-list-photo").First().Children().
		Filter("li.ads-list-photo-item").Filter(":not(.is-adsense):not(.js-booster-inline)").
		Find(".ads-list-photo-item-title").Children().
		Map(func(i int, a *goquery.Selection) string {
			href, _ := a.Attr("href")
			return href
		})
	return adLinks, hasNextPage
}

func getLady(htmlPage *goquery.Document) dto.Lady {
	//TODO: phone MAY BE HIDDEN, handle this case
	phone, _ := htmlPage.Find("dl.js-phone-number").
		Find("a").Attr("href")
	phone = strings.TrimPrefix(phone, "tel:+373")

	return dto.Lady{Phone: phone, UsedLast: time.Now()}
}
