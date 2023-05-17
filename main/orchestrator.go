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

const updateLadiesInterval = time.Hour * 4
const sendOrdersBaseInterval = time.Hour * 1
const sendOrdersIntervalVariation = sendOrdersBaseInterval / 2
const jobStart = time.Hour * 10
const jobEnd = time.Hour * 21

func main() {
	//mainLogic()

	target.TestCookies()
}

func mainLogic() {
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
		sleepAtNight()
		target.OrderItem()

		sleepDuration := time.Duration(float64(sendOrdersIntervalVariation) *
			(rand.Float64() - 0.5))
		time.Sleep(sleepDuration)
	}
}

func sleepAtNight() {
	loc, _ := time.LoadLocation("Local")
	year, month, day := time.Now().In(loc).Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, loc)

	currentTime := time.Now()
	startTime := midnight.Add(jobStart)
	endTime := midnight.Add(jobEnd)

	if currentTime.Before(startTime) {
		sleepDuration := startTime.Sub(currentTime)
		time.Sleep(sleepDuration)
	} else if currentTime.After(endTime) {
		sleepDuration := startTime.Add(time.Hour * 24).Sub(currentTime)
		time.Sleep(sleepDuration)
	}
}
