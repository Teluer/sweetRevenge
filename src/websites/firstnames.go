package websites

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"sweetRevenge/src/db/dao"
	"sweetRevenge/src/db/dto"
	"sweetRevenge/src/websites/web"
	"sync"
)

func UpdateFirstNames(wg *sync.WaitGroup, firstNamesUrl string) {
	log.Info("Updating first names if needed")
	if dao.Dao.IsTableEmpty(&dto.FirstName{}) {
		log.Info("First names table empty, updating")
		names := fetchFirstNames(firstNamesUrl)
		dao.Dao.Insert(&names)
	}
	wg.Done()
}

func fetchFirstNames(firstNamesUrl string) (dtos []dto.FirstName) {
	page := web.GetUrlUnsafe(firstNamesUrl).Find("td.sur")
	femaleNames := page.Children()

	//getting the most popular names only
	femaleNames.Slice(0, 174).Each(func(_ int, name *goquery.Selection) {
		dtos = append(dtos, dto.FirstName{FirstName: name.Text()})
	})

	return dtos
}
