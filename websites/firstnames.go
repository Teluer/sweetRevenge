package websites

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"sweetRevenge/db/dao"
	"sweetRevenge/db/dto"
	"sweetRevenge/websites/web"
	"sync"
)

const firstNamesUrl = "https://forebears.io/moldova/forenames"

func UpdateFirstNames(wg *sync.WaitGroup) {
	log.Info("Updating first names if needed")
	if dao.IsTableEmpty(&dto.FirstName{}) {
		log.Info("First names table empty, updating")
		names := fetchFirstNames()
		dao.Insert(&names)
	}
	wg.Done()
}

func fetchFirstNames() (dtos []dto.FirstName) {
	page := web.GetUrl(firstNamesUrl, false).Find("tbody")
	femaleNames := page.Find("div.f").Parent().Next().Children()

	femaleNames.Each(func(_ int, name *goquery.Selection) {
		dtos = append(dtos, dto.FirstName{FirstName: name.Text()})
	})

	return dtos
}
