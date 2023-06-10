package websites

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"sweetRevenge/src/db/dao"
	"sweetRevenge/src/db/dto"
	"sweetRevenge/src/websites/web"
	"sync"
)

// UpdateFirstNames populates "first_names" DB table.
// If the table is empty, it fetches first names from the given URL and saves them in DB.
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
	names := web.GetUnsafe(firstNamesUrl).Find("td.sur").Children()

	//getting the most popular names only
	names.Slice(0, 174).Each(func(_ int, name *goquery.Selection) {
		dtos = append(dtos, dto.FirstName{FirstName: name.Text()})
	})

	return dtos
}
