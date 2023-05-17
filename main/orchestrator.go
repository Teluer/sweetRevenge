package main

import (
	"math/rand"
	"sweetRevenge/websites"
	"sweetRevenge/websites/target"
	"sync"
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
const sendManualOrderRefreshInterval = time.Minute

func main() {
	mainLogic()
	//var wg sync.WaitGroup
	//wg.Add(1)
	////target.TestCookies()
	//go target.Server()
	//target.SendTestOrder()
	//wg.Wait()
}

func mainLogic() {
	var wg sync.WaitGroup

	wg.Add(3)
	go websites.UpdateLastNames(&wg)
	go websites.UpdateFirstNames(&wg)
	//wait for the first update to complete, then proceed
	go updateLadiesRoutine(&wg)
	wg.Wait()

	//everything ready, start sending orders
	//TODO: enable this when manually tested ordering and the admin called
	//go sendOrdersRoutine()
	//run a thread allowing to send a custom order manually
	go manualOrdersRoutine()
}

func updateLadiesRoutine(wg *sync.WaitGroup) {
	websites.UpdateLadies()
	wg.Done()
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

func manualOrdersRoutine() {
	for {
		target.ExecuteManualOrder()
		time.Sleep(sendManualOrderRefreshInterval)
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
