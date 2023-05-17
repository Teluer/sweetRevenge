package websites

import (
	"github.com/PuerkitoBio/goquery"
	"sweetRevenge/db/dao"
	"sweetRevenge/db/dto"
	"sweetRevenge/websites/web"
	"time"
)

const firstNamesUrl = "https://forebears.io/moldova/forenames"

func UpdateFirstNames() {
	//TODO: check if table is empty
	if dao.IsTableEmpty(&dto.FirstName{}) {
		names := fetchFirstNames()
		dao.Insert(&names)
	}
}

func fetchFirstNames() (dtos []dto.FirstName) {
	page := web.Fetch(firstNamesUrl, false).Find("tbody")
	femaleNames := page.Find("div.f").Parent().Next().Children()

	femaleNames.Each(func(_ int, name *goquery.Selection) {
		dtos = append(dtos, dto.FirstName{name.Text(), time.Now(), 0})
	})

	return dtos
}