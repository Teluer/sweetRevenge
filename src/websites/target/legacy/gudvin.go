package legacy

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sweetRevenge/src/config"
	"sweetRevenge/src/db/dao"
	"sweetRevenge/src/db/dto"
	"sweetRevenge/src/rabbitmq"
	"sweetRevenge/src/util"
	"sweetRevenge/src/websites/web"
	"time"
)

type OrderSuccess struct {
	Success  int    `json:"success"`
	Redirect string `json:"redirect_location"`
}

// okay this is ugly but convenient
var orderCfg config.OrdersConfig
var manualOrders []*rabbitmq.ManualOrder

func OrderItem(cfg config.OrdersConfig) {
	defer util.RecoverAndLogError("Orders")
	orderCfg = cfg
	//check manually prepared orders, if there are no manual orders then make random order
	if !executeManualOrder() {
		log.Info("Sending random order")
		name, phone := createRandomCustomer()
		orderItemWithCustomer(name, phone)
	}
}

func orderItemWithCustomer(name, phone string) {
	tor := web.OpenAnonymousSession()

	itemId, link := findRandomItem(tor)
	log.Info(fmt.Sprintf("Sending order for (%s, %s, %s) with cookies from %s to %s ",
		name, phone, itemId, link, orderCfg.TargetOrderLink))

	cookies := getCookies(tor, link)
	log.Info("Got cookies: ", cookies)

	req := prepareOrderRequest(orderCfg.TargetOrderLink, name, phone, itemId, link, cookies)
	resp, body := tor.SendRequest(req)
	log.Info(string(body))
	cookies = append(cookies, resp.Cookies()...)

	var responseBody OrderSuccess
	json.Unmarshal(body, &responseBody)

	if responseBody.Success != 1 {
		log.Error("Response body is not success=1!")
		return
	}

	//todo: check referer
	log.Info("Visiting order page to reproduce user behaviour")
	req = prepareConfirmOrderGetRequest(responseBody.Redirect, link, cookies)
	resp, _ = tor.SendRequest(req)
	cookies = append(cookies, resp.Cookies()...)

	log.Info("Confirming payment method for the new order " +
		"(NOT, because the website disabled this button)")
	//req = prepareConfirmOrderRequest(responseBody.Redirect, cookies)
	//resp, body = tor.SendRequest(req)

	saveOrderHistory(name, phone, itemId)
	log.Info("Sent order successfully")
}

func QueueManualOrder(order *rabbitmq.ManualOrder) {
	manualOrders = append(manualOrders, order)
}

func executeManualOrder() bool {
	log.Info("Checking if should send manual orders")
	if len(manualOrders) == 0 {
		log.Info("Manual orders not found")
		return false
	}

	order := manualOrders[0]
	log.Info(fmt.Sprintf("Sending manual order for %s %s", order.Name, order.Phone))
	orderItemWithCustomer(order.Name, order.Phone)
	manualOrders = manualOrders[1:]
	return true
}

func prepareOrderRequest(target, name, phone, itemId, referer string, cookies []*http.Cookie) *http.Request {
	//create request body
	formData := url.Values{}
	formData.Set("variant_id", itemId)
	formData.Set("amount", "")
	formData.Set("IsFastOrder", "true")
	formData.Set("name", name)
	formData.Set("phone", phone)
	formData.Set("delivery_id", "4")

	// Encode the form data
	body := strings.NewReader(formData.Encode())

	//create request
	request, err := http.NewRequest("POST", target, body)
	if err != nil {
		log.WithError(err).Error("Failed to create request for order!")
		panic("failed to create request!")
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
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("Origin", "https://gudvin.md")
	request.Header.Set("DNT", "1")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Referer", referer)
	request.Header.Set("Sec-Fetch-Dest", "empty")
	request.Header.Set("Sec-Fetch-Mode", "cors")
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Host", "gudvin.md")

	return request
}

func prepareConfirmOrderGetRequest(target, referer string, cookies []*http.Cookie) *http.Request {
	request, err := http.NewRequest("GET", target, nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for order!")
		panic("failed to create request!")
	}

	//set cookies previously returned by the server
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	request.Header.Set("Accept-Language", "en-US,en;q=0.5")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("Origin", "https://gudvin.md")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Referer", referer)
	request.Header.Set("Sec-Fetch-Dest", "document")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Set("Sec-Fetch-Mode", "navigate")
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Sec-Fetch-User", "?1")

	return request
}

func prepareConfirmOrderRequest(target string, cookies []*http.Cookie) *http.Request {
	formData := url.Values{}
	formData.Set("payment_method_id", "23")
	formData.Set("checkout", "Применить")
	// Encode the form data
	body := strings.NewReader(formData.Encode())

	request, err := http.NewRequest("POST", target, body)
	if err != nil {
		log.WithError(err).Error("Failed to create request for order!")
		panic("failed to create request!")
	}

	//set cookies previously returned by the server
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	request.Header.Set("Accept-Language", "en-US,en;q=0.5")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("Origin", "https://gudvin.md")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Referer", target)
	request.Header.Set("Sec-Fetch-Dest", "document")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Set("Sec-Fetch-Mode", "navigate")
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Sec-Fetch-User", "?1")

	return request
}

func getCookies(tor *web.AnonymousSession, link string) []*http.Cookie {
	log.Info("Fetching cookies to build order request")
	cookies := tor.FetchCookies(link)
	return cookies
}

func findRandomItem(tor *web.AnonymousSession) (id string, link string) {
	randomCategory := orderCfg.TargetCategories[rand.Intn(len(orderCfg.TargetCategories))]

	log.Info("Fetching random item from category " + randomCategory)

	page := tor.GetUrl(randomCategory)

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

	link = orderCfg.TargetBaselink + link
	log.Info("Will order the following item: " + id + " " + link)
	return id, link
}

func createRandomCustomer() (name string, phone string) {
	const firstNameOnlyIncidence = 0.2
	const firstNameAfterLastNameIncidence = 0.6
	const nameLowerCaseIncidence = 0.05

	log.Info("Generating a random customer name/phone combination")

	//write phones in random formats
	phone = dao.Dao.GetLeastUsedPhone()
	prefixIndex := rand.Intn(len(orderCfg.PhonePrefixes))
	phone = orderCfg.PhonePrefixes[prefixIndex] + phone

	//names should look random as well
	name = dao.Dao.GetLeastUsedFirstName()
	if !evaluateProbability(firstNameOnlyIncidence) {
		if evaluateProbability(firstNameAfterLastNameIncidence) {
			name = dao.Dao.GetLeastUsedLastName() + " " + name
		} else {
			name = name + " " + dao.Dao.GetLeastUsedLastName()
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

func saveOrderHistory(name, phone, itemId string) {
	var record = dto.OrderHistory{
		Phone:         phone,
		Name:          name,
		ItemId:        itemId,
		OrderDateTime: time.Now(),
	}

	dao.Dao.Insert(&record)
}
