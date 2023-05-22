package main

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sweetRevenge/src/config"
	"sweetRevenge/src/db/dao"
	dto2 "sweetRevenge/src/db/dto"
	"sweetRevenge/src/rabbitmq"
	"sweetRevenge/src/websites"
	"sweetRevenge/src/websites/target/legacy"
	"sync"
	"time"
)

func programLogic(cfg config.Config) {
	rand.Seed(time.Now().UnixMilli())
	//wait for the updates to complete, then proceed with orders.
	//this is unnecessary since data integrity checks are in place, keeping this just for lulz
	var wg sync.WaitGroup
	wg.Add(2)
	go websites.UpdateLastNamesRoutine(&wg, cfg.LastNamesUrl)
	go websites.UpdateFirstNamesRoutine(&wg, cfg.FirstNamesUrl)
	wg.Wait()

	go manualOrdersRoutine(cfg.OrdersRoutineCfg.OrdersCfg.Rabbit)
	log.Info("Not STUCK!")

	//TODO: some bug prevents ladies from marking as used
	go updateLadiesRoutine(cfg.LadiesCfg)

	//everything ready, start sending orders
	//go sendOrdersRoutine(cfg.OrdersRoutineCfg)
}

func manualOrdersRoutine(cfg config.RabbitConfig) {
	log.Info("Initializing rabbitmq connection")
	rabbitmq.InitializeRabbitMq(cfg)

	log.Info("Starting manual orders RabbitMq listener")
	for {
		order := rabbitmq.ConsumeManualOrder(cfg.QueueName)
		legacy.QueueManualOrder(order)
		log.Info("Manual order is queued and will be executed by Orders routine")
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

func sendOrdersRoutine(cfg config.OrdersRoutineConfig) {
	log.Info("Starting send orders routine")
	for {
		sleepAtNight(cfg)
		jobStart := time.Now()

		//is everything in place to make orders
		readyToGo := !(dao.IsTableEmpty(&dto2.FirstName{}) || dao.IsTableEmpty(&dto2.LastName{}) || dao.IsTableEmpty(&dto2.Lady{}))
		if readyToGo {
			legacy.OrderItem(cfg.OrdersCfg)
		}

		jobDuration := time.Now().Sub(jobStart)

		sleepDuration := time.Duration(float64(cfg.SendOrdersMaxInterval)*rand.Float64()) - jobDuration
		log.Info("sendOrdersRoutine: sleeping for ", int(sleepDuration/time.Minute), " minutes")
		time.Sleep(sleepDuration)
	}
}

func sleepAtNight(cfg config.OrdersRoutineConfig) {
	loc, _ := time.LoadLocation(cfg.TimeZone)
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
