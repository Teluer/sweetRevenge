package websites

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"sweetRevenge/src/db/dao"
	"sweetRevenge/src/db/dto"
	"sweetRevenge/src/websites/web"
	"sync"
)

func UpdateLastNamesRoutine(wg *sync.WaitGroup, lastNamesUrl string) {
	log.Info("Updating last names if needed")
	if dao.IsTableEmpty(&dto.LastName{}) {
		log.Info("Last names table empty, updating")
		names := fetchLastNames(lastNamesUrl)
		dao.Insert(&names)
	}
	wg.Done()
}

func fetchLastNames(lastNamesUrl string) (dtos []dto.LastName) {
	lastNames := web.GetUrlUnsafe(lastNamesUrl).Find("ol.row").Find("a")
	lastNames.Each(func(_ int, name *goquery.Selection) {
		dtos = append(dtos, dto.LastName{LastName: name.Text()})
	})
	return dtos
}
