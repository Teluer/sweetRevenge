package target

import (
	"fmt"
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

type Order struct {
	OrderCfg           *config.OrdersConfig
	SocksProxy         string
	ConcurrencyCh      chan struct{}
	currentTransaction dao.Database
}

var manualOrders []*rabbitmq.ManualOrder

func (ord *Order) OrderItem() {
	//notify channel when order completed
	defer func() { <-ord.ConcurrencyCh }()

	ord.currentTransaction = dao.Dao.OpenTransaction()
	defer util.RecoverAndRollbackAndLog("Orders", ord.currentTransaction)

	//check manually prepared orders, if there are no manual orders then make random order
	if !ord.executeManualOrder() {
		log.Info("Sending random order")
		name, phone := ord.CreateRandomCustomer()
		ord.orderItemWithCustomer(name, phone)
	}
	ord.currentTransaction.CommitTransaction()
}

func (ord *Order) orderItemWithCustomer(name, phone string) {
	tor := web.OpenAnonymousSession(ord.SocksProxy)
	itemId, link := ord.findRandomItem(tor)

	if ord.OrderCfg.SeleniumEnabled {
		ord.orderItemWithCustomerSelenium(name, phone, itemId, link)
	} else {
		ord.orderItemWithCustomerTor(name, phone, itemId, link, tor)
	}
	ord.saveOrderHistory(name, phone, itemId)
}

func (ord *Order) executeManualOrder() bool {
	log.Info("Checking if should send manual orders")
	if len(manualOrders) == 0 {
		log.Info("Manual orders not found")
		return false
	}

	order := manualOrders[0]
	if order.Name == "" {
		order.Name = ord.generateName()
	}
	if order.Phone == "" {
		order.Phone = ord.generatePhone()
	}

	log.Info(fmt.Sprintf("Sending manual order for %s %s", order.Name, order.Phone))
	ord.orderItemWithCustomer(order.Name, order.Phone)
	manualOrders = manualOrders[1:]
	return true
}

func QueueManualOrder(order *rabbitmq.ManualOrder) {
	manualOrders = append(manualOrders, order)
}

func (ord *Order) CreateRandomCustomer() (name string, phone string) {
	log.Info("Generating a random customer name/phone combination")
	phone = ord.generatePhone()
	name = ord.generateName()
	return
}

func (ord *Order) generateName() string {
	const firstNameOnlyIncidence = 0.2
	const firstNameAfterLastNameIncidence = 0.6
	const nameLowerCaseIncidence = 0.05

	name := ord.currentTransaction.GetLeastUsedFirstName()
	if !ord.evaluateProbability(firstNameOnlyIncidence) {
		if ord.evaluateProbability(firstNameAfterLastNameIncidence) {
			name = ord.currentTransaction.GetLeastUsedLastName() + " " + name
		} else {
			name = name + " " + ord.currentTransaction.GetLeastUsedLastName()
		}
	}
	if ord.evaluateProbability(nameLowerCaseIncidence) {
		name = strings.ToLower(name)
	}
	return name
}

func (ord *Order) generatePhone() string {
	phone := ord.currentTransaction.GetLeastUsedPhone()
	prefixes := strings.Split(ord.OrderCfg.PhonePrefixes, ";")
	prefixIndex := rand.Intn(len(prefixes))
	phone = prefixes[prefixIndex] + phone
	return phone
}

func (ord *Order) evaluateProbability(probability float64) bool {
	return rand.Float64() < probability
}

func (ord *Order) saveOrderHistory(name, phone, itemId string) {
	var record = dto.OrderHistory{
		Phone:         phone,
		Name:          name,
		ItemId:        itemId,
		OrderDateTime: time.Now(),
	}

	ord.currentTransaction.Insert(&record)
}
