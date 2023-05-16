package fetcher

import (
	"github.com/PuerkitoBio/goquery"
	"sweetRevenge/db/dao"
	"sweetRevenge/db/dto"
	"time"
)

const firstNamesUrl = "https://forebears.io/moldova/forenames"
const lastNamesUrl = "https://surnam.es/moldova"

func UpdateFirstNames() {
	//TODO: check if table is empty
	if dao.IsTableEmpty(&dto.FirstName{}) {
		names := fetchFirstNames()
		dao.Insert(&names)
	}
}

func UpdateLastNames() {
	if dao.IsTableEmpty(&dto.LastName{}) {
		names := fetchLastNames()
		dao.Insert(&names)
	}
}

func fetchFirstNames() (dtos []dto.FirstName) {
	page := fetch(firstNamesUrl).Find("tbody")
	femaleNames := page.Find("div.f").Parent().Next().Children()

	femaleNames.Each(func(_ int, name *goquery.Selection) {
		dtos = append(dtos, dto.FirstName{name.Text(), time.Now(), 0})
	})

	return dtos
}

func fetchLastNames() (dtos []dto.LastName) {
	lastNames := fetch(lastNamesUrl).Find("ol.row").Find("a")
	lastNames.Each(func(_ int, name *goquery.Selection) {
		dtos = append(dtos, dto.LastName{name.Text(), time.Now(), 0})
	})
	return dtos
}
