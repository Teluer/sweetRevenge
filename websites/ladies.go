package websites

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"sweetRevenge/db/dao"
	"sweetRevenge/db/dto"
	"sweetRevenge/websites/web"
	"time"
)

const baseUrl = "https://999.md"

// TODO: fetch other sections too
const ladiesUrl = "https://999.md/ru/list/dating-and-greetings/i-need-a-man"
const ladiesPageSleepTime = time.Second * 2
const ladiesAdSleepTime = time.Second * 1

func UpdateLadies() {
	ladies := getLadies()
	dao.SaveNewLadies(ladies)
}

func getLadies() (ladies []dto.Lady) {
	var urls []string
	currentUrl := ladiesUrl
	pageNumber := 1
	for {
		page := web.Fetch(currentUrl, false) // start a goroutine
		ladyUrls, hasNextPage := parseLadiesList(page)
		if len(ladyUrls) > 0 {
			urls = append(urls, ladyUrls...)
		}
		//TODO: remove test condition, make 1 second pause
		if !hasNextPage || true {
			break
		}
		time.Sleep(ladiesPageSleepTime)
		pageNumber++
		currentUrl = ladiesUrl + "?page=" + strconv.Itoa(pageNumber)
	}

	//TODO: use goroutines, make 1 request per second
	for _, url := range urls {
		url = baseUrl + url
		ad := web.Fetch(url, false)
		time.Sleep(ladiesAdSleepTime)
		ladies = append(ladies, getLady(ad))
	}

	//remove duplicated phones
	var uniqueLadies []dto.Lady
MAIN_LOOP:
	for _, lady := range ladies {
		for _, resultLady := range uniqueLadies {
			if resultLady.Phone == lady.Phone {
				continue MAIN_LOOP
			}
		}
		uniqueLadies = append(uniqueLadies, lady)
	}

	return uniqueLadies
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
