package target

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sweetRevenge/db/dao"
	"sweetRevenge/db/dto"
	"sweetRevenge/websites/web"
)

type OrderBody struct {
	//variant_id=158
	//amount
	//IsFastOrder=true
	//name=Сергей Реулец
	//phone=079744327
	//delivery_id=4
	VariantId   string `json:"variant_id"`
	Amount      string `json:"amount"`
	IsFastOrder string `json:"IsFastOrder"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	DeliveryId  string `json:"delivery_id"`
}

// add 0 several times to increase probability
var phonePrefixes = []string{
	"0", "0", "0", "+373", "(+373) ", "+373 ",
}

const baselink = "https://gudvin.md"
const orderLink = baselink + "/okay-cms/fast-order/create-order"

var categories = []string{
	"https://gudvin.md/catalog/ulichnoe-osveschenie",
	"https://gudvin.md/catalog/tovary-dlya-avto",
	"https://gudvin.md/catalog/prochie-tovary",
	"https://gudvin.md/catalog/stereosistemyusiliteli",
	"https://gudvin.md/catalog/melkaya-bytovaya-tehnika",
	"https://gudvin.md/catalog/tovary-dlya-kuhni",
	"https://gudvin.md/catalog/turizm-sport-i-otdyh",
}

func OrderItem() {
	name, phone := createRandomCustomer()
	OrderItemWithCustomer(name, phone)
}

func OrderItemWithCustomer(name, phone string) {
	OrderItemWithCustomerAndTarget(orderLink, name, phone)
}

func OrderItemWithCustomerAndTarget(targetUrl, name, phone string) {
	//TODO: remove println
	fmt.Println(name, phone)
	itemId, link := fetchRandomItem()
	OrderItemWithCustomerAndTargetAndItemAndLink(targetUrl, name, phone, itemId, link)
}

func OrderItemWithCustomerAndTargetAndItemAndLink(targetUrl, name, phone, itemId, link string) {
	fmt.Println(itemId, link)

	cookies := getCookies(link)
	req := prepareOrderPostRequest(targetUrl, name, phone, itemId, link, cookies)
	web.Post(req, true)
}

func ExecuteManualOrder() {
	var manualOrder dto.ManualOrder
	dao.FindFirstAndDelete(&manualOrder)

	//if not found
	if manualOrder.Phone == "" {
		return
	}

	//send either to default target, or to custom url
	if manualOrder.Target == "" {
		OrderItemWithCustomer(manualOrder.Name, manualOrder.Phone)
	} else {
		OrderItemWithCustomerAndTarget(manualOrder.Target, manualOrder.Name, manualOrder.Phone)
	}
}

func prepareOrderPostRequest(target, name, phone, itemId, referer string, cookies []*http.Cookie) *http.Request {
	//create request body
	order, err := json.Marshal(OrderBody{
		VariantId:   itemId,
		Amount:      "",
		Name:        name,
		Phone:       phone,
		DeliveryId:  "4",
		IsFastOrder: "true",
	})
	if err != nil {
		panic("failed to marshal request body!")
	}
	orderBody := string(order)

	//create request
	request, err := http.NewRequest("Post", target, strings.NewReader(orderBody))
	if err != nil {
		panic("failed to make request!")
	}

	//set cookies previously returned by the server
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	//set proper headers
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0")
	request.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	request.Header.Set("Accept-Language", "en-US,en;q=0.5")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("Content-Length", strconv.Itoa(len(orderBody)))
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("Origin", "https://gudvin.md")
	request.Header.Set("DNT", "1")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Referer", referer)
	request.Header.Set("Sec-Fetch-Dest", "empty")
	request.Header.Set("Sec-Fetch-Mode", "cors")
	request.Header.Set("Sec-Fetch-Site", "same-origin")

	return request
}

func getCookies(link string) []*http.Cookie {
	_, cookies := web.FetchWithCookies(link, true)
	return cookies
}

func fetchRandomItem() (id string, link string) {
	randomCategory := categories[rand.Intn(len(categories))]
	page := web.Fetch(randomCategory, true)

	items := page.Find("a.product_preview__name_link")
	randomItem := rand.Intn(items.Length())
	items.EachWithBreak(func(i int, item *goquery.Selection) bool {
		if i == randomItem {
			id, _ = item.Attr("data-product")
			link, _ = item.Attr("href")
			return false
		}
		return true
	})

	link = baselink + link

	return id, link
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
