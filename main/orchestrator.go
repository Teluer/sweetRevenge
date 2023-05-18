package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sweetRevenge/config"
	"sweetRevenge/websites"
	"sweetRevenge/websites/target"
	"sync"
	"time"
)

const sendOrdersBaseInterval = time.Hour * 1
const sendOrdersIntervalVariation = sendOrdersBaseInterval / 2
const sendManualOrderRefreshInterval = time.Minute

func programLogic(cfg config.Config) {
	rand.Seed(time.Now().UnixMilli())

	var wg sync.WaitGroup
	wg.Add(2)
	go websites.UpdateLastNamesRoutine(&wg, cfg.LastNamesUrl)
	go websites.UpdateFirstNamesRoutine(&wg, cfg.FirstNamesUrl)
	//wait for the first update to complete, then proceed
	go updateLadiesRoutine(&wg, cfg.LadiesCfg)
	wg.Wait()

	//everything ready, start sending orders
	//TODO: enable this when manually tested ordering and the operator called
	//go sendOrdersRoutine(cfg.OrdersRoutineCfg)
	//run a thread allowing to send a custom order manually
	go manualOrdersRoutine()
}

func updateLadiesRoutine(wg *sync.WaitGroup, cfg config.LadiesConfig) {
	log.Info("Starting update ladies routine")
	websites.UpdateLadies(cfg.LadiesBaseUrl, cfg.LadiesUrls)
	wg.Done()
	for {
		log.Info(fmt.Sprintf("updateLadiesRoutine: sleeping for %d minutes",
			cfg.UpdateLadiesInterval/time.Minute))
		time.Sleep(cfg.UpdateLadiesInterval)
		websites.UpdateLadies(cfg.LadiesBaseUrl, cfg.LadiesUrls)
	}
}

func sendOrdersRoutine(cfg config.OrdersRoutineConfig) {
	log.Info("Starting send orders routine")
	for {
		sleepAtNight(cfg)
		target.OrderItem(cfg.OrdersCfg)

		//TODO: make variation relative
		sleepDuration := cfg.SendOrdersBaseInterval +
			time.Duration(float64(sendOrdersIntervalVariation)*(rand.Float64()-0.5))
		log.Info(fmt.Sprintf("sendOrdersRoutine: sleeping for %d minutes",
			sleepDuration/time.Minute))
		time.Sleep(sleepDuration)
	}
}

// TODO: do this inside orders routine to keep normal order rates
func manualOrdersRoutine() {
	log.Info("Starting manual orders routine")
	for {
		target.ExecuteManualOrder()
		log.Info(fmt.Sprintf("manualOrdersRoutine: sleeping for %d minutes",
			sendManualOrderRefreshInterval/time.Minute))
		time.Sleep(sendManualOrderRefreshInterval)
	}
}

// TODO: test this
func sleepAtNight(cfg config.OrdersRoutineConfig) {
	loc, _ := time.LoadLocation("Local")
	year, month, day := time.Now().In(loc).Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, loc)

	currentTime := time.Now()
	startTime := midnight.Add(cfg.DayStart)
	endTime := midnight.Add(cfg.DayEnd)

	if currentTime.Before(startTime) {
		sleepDuration := startTime.Sub(currentTime)
		log.Info("Beyond work hours, sleeping until " +
			time.Now().Add(sleepDuration).Format("2006-01-02 15:04:05"))
		time.Sleep(sleepDuration)
	} else if currentTime.After(endTime) {
		sleepDuration := startTime.Add(time.Hour * 24).Sub(currentTime)
		log.Info("Beyond work hours, sleeping until " +
			time.Now().Add(sleepDuration).Format("2006-01-02 15:04:05"))
		time.Sleep(sleepDuration)
	}
}
