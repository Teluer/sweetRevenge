package fetcher

import (
	"github.com/PuerkitoBio/goquery"
	"sweetRevenge/fetcher/web"
)

var categories = []string{
	"https://gudvin.md/catalog/ulichnoe-osveschenie",
	"https://gudvin.md/catalog/tovary-dlya-avto",
	"https://gudvin.md/catalog/prochie-tovary",
	"https://gudvin.md/catalog/stereosistemyusiliteli",
	"https://gudvin.md/catalog/melkaya-bytovaya-tehnika",
	"https://gudvin.md/catalog/tovary-dlya-kuhni",
	"https://gudvin.md/catalog/turizm-sport-i-otdyh",
}

type ItemDto struct {
	id       string
	link     string
	category string
}

// TODO: this fetches relative urls for each product
func FetchGoods() []ItemDto {
	var itemDtos []ItemDto
	for _, category := range categories {
		items := web.Browser{}.Fetch(category, false).Find("a.product_preview__name_link")
		items.Each(func(_ int, item *goquery.Selection) {
			id, _ := item.Attr("data-product")
			link, _ := item.Attr("href")
			itemDtos = append(itemDtos, ItemDto{
				id,
				link,
				category,
			})
		})
	}
	return itemDtos
}
