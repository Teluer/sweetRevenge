package fetcher

import (
	"github.com/PuerkitoBio/goquery"
	"sweetRevenge/db/dao"
	"sweetRevenge/db/dto"
	"sweetRevenge/web"
	"time"
)

const lastNamesUrl = "https://surnam.es/moldova"

func UpdateLastNames() {
	if dao.IsTableEmpty(&dto.LastName{}) {
		names := fetchLastNames()
		dao.Insert(&names)
	}
}

func fetchLastNames() (dtos []dto.LastName) {
	lastNames := web.Fetch(lastNamesUrl, false).Find("ol.row").Find("a")
	lastNames.Each(func(_ int, name *goquery.Selection) {
		dtos = append(dtos, dto.LastName{name.Text(), time.Now(), 0})
	})
	return dtos
}
