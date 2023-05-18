package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"os"
	"sweetRevenge/websites"
	"sweetRevenge/websites/target"
	"sync"
	"time"
)

const updateLadiesInterval = time.Hour * 4
const sendOrdersBaseInterval = time.Hour * 1
const sendOrdersIntervalVariation = sendOrdersBaseInterval / 2
const jobStart = time.Hour * 10
const jobEnd = time.Hour * 21
const sendManualOrderRefreshInterval = time.Minute

func init() {
	//log.SetReportCaller(true)
	log.Info("Program Startup")
	//TODO: write to file as well
	log.SetOutput(io.MultiWriter(os.Stdout))
}

func main() {

	mainLogic()
	//test.TestAnonSending()
	//test.SendTestRequest()
	//websites.UpdateLadies()

	//wait indefinitely
	select {}
}

func mainLogic() {
	rand.Seed(time.Now().UnixMilli())

	log.Info("Updating first and last names if needed")
	var wg sync.WaitGroup
	wg.Add(2)
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
	log.Info("Starting update ladies routine")
	websites.UpdateLadies()
	wg.Done()
	for {
		time.Sleep(updateLadiesInterval)
		websites.UpdateLadies()
	}
}

func sendOrdersRoutine() {
	log.Info("Starting send orders routine")
	for {
		sleepAtNight()
		target.OrderItem()

		sleepDuration := time.Duration(float64(sendOrdersIntervalVariation) *
			(rand.Float64() - 0.5))
		time.Sleep(sleepDuration)
	}
}

func manualOrdersRoutine() {
	log.Info("Starting manual orders routine")
	for {
		target.ExecuteManualOrder()
		time.Sleep(sendManualOrderRefreshInterval)
	}
}

// TODO: test this
func sleepAtNight() {
	loc, _ := time.LoadLocation("Local")
	year, month, day := time.Now().In(loc).Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, loc)

	currentTime := time.Now()
	startTime := midnight.Add(jobStart)
	endTime := midnight.Add(jobEnd)

	if currentTime.Before(startTime) {
		log.Info("Beyond work hours, sleeping")
		sleepDuration := startTime.Sub(currentTime)
		time.Sleep(sleepDuration)
	} else if currentTime.After(endTime) {
		log.Info("Beyond work hours, sleeping")
		sleepDuration := startTime.Add(time.Hour * 24).Sub(currentTime)
		time.Sleep(sleepDuration)
	}
}
