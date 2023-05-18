package websites

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"sweetRevenge/db/dao"
	"sweetRevenge/db/dto"
	"sweetRevenge/websites/web"
	"sync"
)

const lastNamesUrl = "https://surnam.es/moldova"

func UpdateLastNames(wg *sync.WaitGroup) {
	log.Info("Updating last names if needed")
	if dao.IsTableEmpty(&dto.LastName{}) {
		log.Info("Last names table empty, updating")
		names := fetchLastNames()
		dao.Insert(&names)
	}
	wg.Done()
}

func fetchLastNames() (dtos []dto.LastName) {
	lastNames := web.GetUrl(lastNamesUrl, false).Find("ol.row").Find("a")
	lastNames.Each(func(_ int, name *goquery.Selection) {
		dtos = append(dtos, dto.LastName{LastName: name.Text()})
	})
	return dtos
}
