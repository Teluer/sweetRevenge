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

var orders struct {
	orderCfg           config.OrdersConfig
	manualOrders       []*rabbitmq.ManualOrder
	currentTransaction dao.Database
}

func OrderItem(cfg config.OrdersConfig, socksProxy string) {
	orders.orderCfg = cfg
	orders.currentTransaction = dao.Dao.OpenTransaction()
	defer util.RecoverAndRollbackAndLog("Orders", orders.currentTransaction)

	//check manually prepared orders, if there are no manual orders then make random order
	if !executeManualOrder(socksProxy) {
		log.Info("Sending random order")
		name, phone := CreateRandomCustomer()
		orderItemWithCustomer(name, phone, socksProxy)
	}
	orders.currentTransaction.CommitTransaction()
	orders.currentTransaction = nil
}

func orderItemWithCustomer(name, phone, socksProxy string) {
	tor := web.OpenAnonymousSession(socksProxy)
	itemId, link := findRandomItem(tor)

	if orders.orderCfg.SeleniumEnabled {
		orderItemWithCustomerSelenium(name, phone, itemId, link, socksProxy)
	} else {
		orderItemWithCustomerTor(name, phone, itemId, link, tor)
	}
	saveOrderHistory(name, phone, itemId)
}

func executeManualOrder(socksProxy string) bool {
	log.Info("Checking if should send manual orders")
	if len(orders.manualOrders) == 0 {
		log.Info("Manual orders not found")
		return false
	}

	order := orders.manualOrders[0]
	if order.Name == "" {
		order.Name = generateName()
	}
	if order.Phone == "" {
		order.Phone = generatePhone()
	}

	log.Info(fmt.Sprintf("Sending manual order for %s %s", order.Name, order.Phone))
	orderItemWithCustomer(order.Name, order.Phone, socksProxy)
	orders.manualOrders = orders.manualOrders[1:]
	return true
}

func QueueManualOrder(order *rabbitmq.ManualOrder) {
	orders.manualOrders = append(orders.manualOrders, order)
}

func CreateRandomCustomer() (name string, phone string) {
	log.Info("Generating a random customer name/phone combination")
	phone = generatePhone()
	name = generateName()
	return
}

func generateName() string {
	const firstNameOnlyIncidence = 0.2
	const firstNameAfterLastNameIncidence = 0.6
	const nameLowerCaseIncidence = 0.05

	name := orders.currentTransaction.GetLeastUsedFirstName()
	if !evaluateProbability(firstNameOnlyIncidence) {
		if evaluateProbability(firstNameAfterLastNameIncidence) {
			name = orders.currentTransaction.GetLeastUsedLastName() + " " + name
		} else {
			name = name + " " + orders.currentTransaction.GetLeastUsedLastName()
		}
	}
	if evaluateProbability(nameLowerCaseIncidence) {
		name = strings.ToLower(name)
	}
	return name
}

func generatePhone() string {
	phone := orders.currentTransaction.GetLeastUsedPhone()
	prefixes := strings.Split(orders.orderCfg.PhonePrefixes, ";")
	prefixIndex := rand.Intn(len(prefixes))
	phone = prefixes[prefixIndex] + phone
	return phone
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

	orders.currentTransaction.Insert(&record)
}
