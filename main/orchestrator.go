package main

import (
	"math/rand"
	"sweetRevenge/websites"
	"sweetRevenge/websites/target"
	"time"
)

/*TODO:
connect to TOR via Socks
go to the shop and make an order
*/

var updateLadiesInterval = time.Hour * 4
var sendOrdersBaseInterval = time.Hour * 1
var sendOrdersIntervalVariation = sendOrdersBaseInterval / 2

func main() {
	websites.UpdateLastNames()
	websites.UpdateFirstNames()
	//wait for the first update to complete, then start goroutine
	websites.UpdateLadies()
	go updateLadiesRoutine()
	//everything ready, start sending orders
	go sendOrdersRoutine()
}

func updateLadiesRoutine() {
	for {
		time.Sleep(updateLadiesInterval)
		websites.UpdateLadies()
	}
}

func sendOrdersRoutine() {
	for {
		target.OrderItem()

		sleepDuration := time.Duration(float64(sendOrdersIntervalVariation) *
			(rand.Float64() - 0.5))
		time.Sleep(sleepDuration)
	}
}
