package fetcher

import (
	"github.com/PuerkitoBio/goquery"
	"sweetRevenge/web"
	"time"
)

type Item struct {
	id   string
	link string
}

var categories = []string{
	"https://gudvin.md/catalog/ulichnoe-osveschenie",
	"https://gudvin.md/catalog/tovary-dlya-avto",
	"https://gudvin.md/catalog/prochie-tovary",
	"https://gudvin.md/catalog/stereosistemyusiliteli",
	"https://gudvin.md/catalog/melkaya-bytovaya-tehnika",
	"https://gudvin.md/catalog/tovary-dlya-kuhni",
	"https://gudvin.md/catalog/turizm-sport-i-otdyh",
}

// TODO: fetch a random item from a random category
func OrderItem() {

}

func OrderManyItems(amount int, delay time.Duration) {
	for i := amount; i > 0; i-- {
		OrderItem()
		time.Sleep(delay)
	}
}

// TODO: this fetches relative urls for each product
func fetchItems() []Item {
	var itemDtos []Item
	for _, category := range categories {
		items := web.Fetch(category, false).Find("a.product_preview__name_link")
		items.Each(func(_ int, item *goquery.Selection) {
			id, _ := item.Attr("data-product")
			link, _ := item.Attr("href")
			itemDtos = append(itemDtos, Item{
				id:   id,
				link: link,
			})
		})
	}
	return itemDtos
}
