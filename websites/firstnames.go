package websites

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"sweetRevenge/db/dao"
	"sweetRevenge/db/dto"
	"sweetRevenge/websites/web"
	"sync"
)

func UpdateFirstNamesRoutine(wg *sync.WaitGroup, firstNamesUrl string) {
	log.Info("Updating first names if needed")
	if dao.IsTableEmpty(&dto.FirstName{}) {
		log.Info("First names table empty, updating")
		names := fetchFirstNames(firstNamesUrl)
		dao.Insert(&names)
	}
	wg.Done()
}

func fetchFirstNames(firstNamesUrl string) (dtos []dto.FirstName) {
	page := web.GetUrl(firstNamesUrl, false).Find("tbody")
	femaleNames := page.Find("div.f").Parent().Next().Children()

	femaleNames.Each(func(_ int, name *goquery.Selection) {
		dtos = append(dtos, dto.FirstName{FirstName: name.Text()})
	})

	return dtos
}
