package websites

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"sweetRevenge/src/db/dao"
	"sweetRevenge/src/db/dto"
	"sweetRevenge/src/util"
	"sweetRevenge/src/websites/web"
	"time"
)

func UpdateLadies(ladiesBaseUrl string, ladiesUrls []string) {
	defer util.RecoverAndLogError("LadiesUpdate")

	log.Info("Ladies update triggered")
	ladies := getLadies(ladiesBaseUrl, ladiesUrls)
	log.Info(fmt.Sprintf("Found %d ladies", len(ladies)))
	dao.SaveNewLadies(ladies)
}

func getLadies(ladiesBaseUrl string, ladiesUrls []string) (ladies []dto.Lady) {
	var urls []string
	for _, ladyCategory := range ladiesUrls {
		for currentUrl, pageNumber := ladyCategory, 1; ; currentUrl = ladyCategory + "?page=" + strconv.Itoa(pageNumber) {
			log.Info("Fetching lady list from " + currentUrl)
			page := web.GetUrl(currentUrl, false)
			ladyUrls, hasNextPage := parseLadiesList(page)
			if len(ladyUrls) > 0 {
				urls = append(urls, ladyUrls...)
			}
			if !hasNextPage {
				break
			}
			pageNumber++
			currentUrl = ladyCategory + "?page=" + strconv.Itoa(pageNumber)
		}
	}

	//send all requests consecutively to avoid getting blocked
	log.Info("Fetching lady phones from all categories")
	for _, url := range urls {
		url = ladiesBaseUrl + url
		request := getRequestWithPopupBypass(url)
		ad := web.GetRequest(request, true)
		lady := getLady(ad)
		if lady.Phone != "" {
			ladies = append(ladies, lady)
		}
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
	phone, _ := htmlPage.Find("dl.js-phone-number").
		Find("a").Attr("href")
	phone = strings.TrimPrefix(phone, "tel:+373")

	return dto.Lady{Phone: phone}
}

func getRequestWithPopupBypass(url string) *http.Request {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request to get lady!")
		return nil
	}

	request.AddCookie(&http.Cookie{
		Name:    "age_popup_show_guest",
		Value:   "False",
		Expires: time.Now().Add(365 * 24 * time.Hour),
	})
	request.AddCookie(&http.Cookie{
		Name:    "age_popup_show",
		Value:   "False",
		Expires: time.Now().Add(365 * 24 * time.Hour),
	})

	return request
}
