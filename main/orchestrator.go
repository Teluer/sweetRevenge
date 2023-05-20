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

func programLogic(cfg config.Config) {
	rand.Seed(time.Now().UnixMilli())

	var wg sync.WaitGroup
	wg.Add(3)
	go websites.UpdateLastNamesRoutine(&wg, cfg.LastNamesUrl)
	go websites.UpdateFirstNamesRoutine(&wg, cfg.FirstNamesUrl)
	//wait for the first update to complete, then proceed
	go updateLadiesRoutine(&wg, cfg.LadiesCfg)
	wg.Wait()

	//everything ready, start sending orders
	go sendOrdersRoutine(cfg.OrdersRoutineCfg)
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

		variationCoefficient := cfg.SendOrdersIntervalVariation*(rand.Float64()-0.5) + 1
		sleepDuration := time.Duration(float64(cfg.SendOrdersBaseInterval) * variationCoefficient)
		log.Info(fmt.Sprintf("sendOrdersRoutine: sleeping for %d minutes",
			sleepDuration/time.Minute))
		time.Sleep(sleepDuration)
	}
}

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
