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

func UpdateLadies(ladiesBaseUrl string, ladiesUrls []string, socksProxy string) {
	defer util.RecoverAndLog("LadiesUpdate")

	log.Info("Ladies update triggered")
	totalInserted := 0
	tor := web.OpenAnonymousSession(socksProxy)
	for _, ladyCategory := range ladiesUrls {
		totalInserted += fetchLadies(ladiesBaseUrl, ladyCategory, tor)
	}
	log.Info(fmt.Sprintf("Found %d new ladies", totalInserted))
}

func fetchLadies(ladiesBaseUrl string, ladyCategory string, tor *web.AnonymousSession) (insertedCount int) {
	var urls []string
	log.Info("Fetching lady list from " + ladyCategory)

	//the first time I actually needed a do/while loop in my life
	for pageNumber, hasNext := 1, true; hasNext; pageNumber++ {
		_, page := tor.GetUrl(ladyCategory + "?page=" + strconv.Itoa(pageNumber))
		ladyUrls := parseLadiesList(page)
		if len(ladyUrls) > 0 {
			urls = append(urls, ladyUrls...)
		}
		hasNext = hasNextPage(page)
	}

	//send all requests consecutively to avoid getting blocked
	log.Info("Fetching lady phones from " + ladyCategory)
	var ladies []dto.Lady
	for _, url := range urls {
		url = ladiesBaseUrl + url
		request := prepareRequestWithPopupBypass(url)
		if request == nil {
			continue
		}
		_, ad := tor.GetRequest(request)
		lady := getLady(ad)
		if lady.Phone != "" {
			ladies = append(ladies, lady)
		}
	}
	ladies = removeDuplicates(ladies)
	return dao.Dao.SaveNewLadies(ladies)
}

func parseLadiesList(htmlPage *goquery.Document) (adLinks []string) {
	adLinks = htmlPage.Find("ul.ads-list-photo").First().Children().
		Filter("li.ads-list-photo-item").
		Filter(":not(.is-adsense):not(.js-booster-inline)").
		Find(".ads-list-photo-item-title").Children().
		Map(func(i int, a *goquery.Selection) string {
			href, _ := a.Attr("href")
			return href
		})
	return adLinks
}

func hasNextPage(htmlPage *goquery.Document) bool {
	return htmlPage.Find("nav.paginator").
		Find(".current").Next().Length() > 0
}

func getLady(htmlPage *goquery.Document) dto.Lady {
	phone, _ := htmlPage.Find("dl.js-phone-number").
		Find("a").Attr("href")

	if !strings.HasPrefix(phone, "tel:+373") {
		log.Warn("Bad phone prefix in a lady: ", phone)
		return dto.Lady{}
	}
	phone = strings.TrimPrefix(phone, "tel:+373")

	return dto.Lady{Phone: phone}
}

func prepareRequestWithPopupBypass(url string) *http.Request {
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

func removeDuplicates(ladies []dto.Lady) (uniqueLadies []dto.Lady) {
MAIN_LOOP:
	for _, lady := range ladies {
		for _, resultLady := range uniqueLadies {
			if resultLady.Phone == lady.Phone {
				continue MAIN_LOOP
			}
		}
		uniqueLadies = append(uniqueLadies, lady)
	}
	return
}
