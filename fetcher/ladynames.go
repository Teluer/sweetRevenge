package fetcher

import (
	"github.com/PuerkitoBio/goquery"
)

const firstNamesUrl = "https://forebears.io/moldova/forenames"
const lastNamesUrl = "https://surnam.es/moldova"

func FetchFirstNames() []string {
	page := fetch(firstNamesUrl).Find("tbody")
	females := page.Find("div.f")
	femaleNames := females.Parent().Next().Children()
	femaleNamesString := femaleNames.Map(func(_ int, name *goquery.Selection) string {
		return name.Text()
	})
	return femaleNamesString
}

func FetchLastNames() []string {
	page := fetch(firstNamesUrl).Find("ol.row")
	lastNames := page.Find("a")
	lastNamesString := lastNames.Map(func(_ int, name *goquery.Selection) string {
		return name.Text()
	})
	return lastNamesString
}
