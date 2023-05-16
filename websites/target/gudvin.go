package target

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math/rand"
	"strings"
	"sweetRevenge/db/dao"
	"sweetRevenge/websites/web"
	"time"
)

type Item struct {
	id   string
	link string
}

// add 0 several times to increase probability
var phonePrefixes = []string{
	"0", "0", "0", "+373", "(+373) ", "+373 ",
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

// TODO: fetch a random item from a random category, make order
func OrderItem() {
	name, phone := createRandomCustomer()
	//TODO: remove println
	fmt.Println(name, phone)

	// TODO: fetch a random item from a random category, make order
}

func OrderManyItems(amount int, delay time.Duration) {
	for i := amount; i > 0; i-- {
		OrderItem()
		time.Sleep(delay)
	}
}

// TODO: this fetches relative urls for each product
// TODO: get random item from random category
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

func createRandomCustomer() (name string, phone string) {
	const firstNameOnlyIncidence = 0.25
	const firstNameAfterLastNameIncidence = 0.6
	const nameLowerCaseIncidence = 0.08
	const phoneWithSpaceIncidence = 0.5

	//write phones in random formats
	phone = dao.GetLeastUsedPhone()
	prefixIndex := rand.Intn(len(phonePrefixes))
	if prefixIndex >= 4 && evaluateProbability(phoneWithSpaceIncidence) {
		phone = phone[:2] + " " + phone[2:]
	}
	prefix := phonePrefixes[rand.Intn(len(phonePrefixes))]
	phone = prefix + phone

	//names should look random as well
	name = dao.GetLeastUsedFirstName()
	if !evaluateProbability(firstNameOnlyIncidence) {
		if evaluateProbability(firstNameAfterLastNameIncidence) {
			name = dao.GetLeastUsedLastName() + " " + name
		} else {
			name = name + " " + dao.GetLeastUsedLastName()
		}
	}
	if evaluateProbability(nameLowerCaseIncidence) {
		name = strings.ToLower(name)
	}

	return name, phone
}

func evaluateProbability(probability float64) bool {
	return rand.Float64() < probability
}
