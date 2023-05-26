package main

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sweetRevenge/src/admin"
	"sweetRevenge/src/config"
	"sweetRevenge/src/db/dao"
	dto2 "sweetRevenge/src/db/dto"
	"sweetRevenge/src/rabbitmq"
	"sweetRevenge/src/util"
	"sweetRevenge/src/websites"
	"sweetRevenge/src/websites/target/legacy"
	"sync"
	"time"
)

func programLogic(cfg *config.Config) {
	rand.Seed(time.Now().UnixMilli())
	//wait for the updates to complete, then proceed with orders.
	//this is unnecessary since data integrity checks are in place, keeping this just for lulz
	var wg sync.WaitGroup
	wg.Add(2)
	go websites.UpdateLastNamesRoutine(&wg, cfg.LastNamesUrl)
	go websites.UpdateFirstNamesRoutine(&wg, cfg.FirstNamesUrl)
	wg.Wait()

	go manualOrdersRoutine()
	go updateLadiesRoutine(cfg.LadiesCfg)
	//everything ready, start sending orders
	go sendOrdersRoutine(&cfg.OrdersRoutineCfg)

	go admin.ControlPanel(&cfg.OrdersRoutineCfg)
}

func manualOrdersRoutine() {
	log.Info("Starting manual orders RabbitMq listener")
	for {
		func() {
			defer util.RecoverAndLogError("RabbitMq")
			order := rabbitmq.ConsumeManualOrder()
			legacy.QueueManualOrder(order)
			log.Info("Manual order is queued and will be executed by Orders routine")
		}()
	}
}

func updateLadiesRoutine(cfg config.LadiesConfig) {
	log.Info("Starting update ladies routine")
	for {
		websites.UpdateLadies(cfg.LadiesBaseUrl, cfg.LadiesUrls)
		log.Info("updateLadiesRoutine: sleeping for ", int(cfg.UpdateLadiesInterval/time.Minute), " minutes")
		time.Sleep(cfg.UpdateLadiesInterval)
	}
}

func sendOrdersRoutine(cfg *config.OrdersRoutineConfig) {
	log.Info("Starting send orders routine")

	//sleeping at first to avoid order spamming due to multiple restarts
	sleepDuration := time.Duration(float64(cfg.SendOrdersMaxInterval) * rand.Float64())
	log.Infof("sendOrdersRoutine: sending order in %.2f minutes", float64(sleepDuration/time.Minute))
	time.Sleep(sleepDuration)

	for {
		log.Info("sendOrdersRoutine: Order flow triggered")
		sleepAtNight(cfg)

		//is everything in place to make orders
		readyToGo := !(dao.Dao.IsTableEmpty(&dto2.FirstName{}) ||
			dao.Dao.IsTableEmpty(&dto2.LastName{}) ||
			dao.Dao.IsTableEmpty(&dto2.Lady{}))
		ordersEnabled := cfg.SendOrdersEnabled

		if readyToGo && ordersEnabled {
			go legacy.OrderItem(cfg.OrdersCfg)
		} else {
			if !readyToGo {
				log.Warn("Cannot send orders due to empty database tables, please check DB!")
			}
			if !ordersEnabled {
				log.Info("SendOrdersEnabled = false, not sending anything")
			}
		}

		sleepDuration := time.Duration(float64(cfg.SendOrdersMaxInterval) * rand.Float64())
		log.Infof("sendOrdersRoutine: scheduling next order in %.2f minutes", float64(sleepDuration)/float64(time.Minute))
		time.Sleep(sleepDuration)
	}
}

func sleepAtNight(cfg *config.OrdersRoutineConfig) {
	loc, _ := time.LoadLocation(cfg.TimeZone)
	year, month, day := time.Now().In(loc).Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, loc)

	currentTime := time.Now()
	startTime := midnight.Add(cfg.DayStart)
	endTime := midnight.Add(cfg.DayEnd)

	var sleepDuration time.Duration
	if currentTime.Before(startTime) {
		sleepDuration = startTime.Sub(currentTime)
	} else if currentTime.After(endTime) {
		sleepDuration = startTime.Add(time.Hour * 24).Sub(currentTime)
	} else {
		return
	}
	log.Info("sendOrdersRoutine: Beyond work hours, sleeping until " +
		time.Now().Add(sleepDuration).Format("2006-01-02 15:04:05"))
	time.Sleep(sleepDuration)
}
