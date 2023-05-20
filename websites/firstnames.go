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
	page := web.GetUnsafe(firstNamesUrl).Find("td.sur")
	femaleNames := page.Children()

	//getting the most popular names only
	femaleNames.Slice(0, 174).Each(func(_ int, name *goquery.Selection) {
		dtos = append(dtos, dto.FirstName{FirstName: name.Text()})
	})

	return dtos
}
