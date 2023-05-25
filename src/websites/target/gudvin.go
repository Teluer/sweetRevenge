package target

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"math/rand"
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
	PanicIfVpnNotEnabled()

	itemId, link := findRandomItem()
	log.Info(fmt.Sprintf("Sending order for (%s, %s, %s)",
		name, phone, itemId))

	selenium := Connect(link)
	defer selenium.Close()

	selenium.Click("a.fn_fast_order_button")
	time.Sleep(time.Second * 3)
	selenium.Input("input.fn_validate_fast_name", name)
	selenium.Input("input.fn_validate_fast_phone", phone)
	selenium.Click("input.fn_fast_order_submit")

	selenium.SolveCaptcha()
	//todo: click confirm payment method button

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

func findRandomItem() (id string, link string) {
	randomCategory := orderCfg.TargetCategories[rand.Intn(len(orderCfg.TargetCategories))]

	log.Info("Fetching random item from category " + randomCategory)

	page := web.GetUrlUnsafe(randomCategory)

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
