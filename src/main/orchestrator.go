package main

import (
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sweetRevenge/src/admin"
	"sweetRevenge/src/config"
	"sweetRevenge/src/db/dao"
	"sweetRevenge/src/rabbitmq"
	"sweetRevenge/src/util"
	"sweetRevenge/src/websites"
	"sweetRevenge/src/websites/target"
	"sync"
	"time"
)

func programLogic(cfg *config.Config) {
	loc, _ := time.LoadLocation(cfg.TimeZone)

	//wait for the updates to complete, then proceed with orders.
	//this is unnecessary since data integrity checks are in place, keeping this just for lulz
	var wg sync.WaitGroup
	wg.Add(2)
	go websites.UpdateLastNames(&wg, cfg.LastNamesUrl)
	go websites.UpdateFirstNames(&wg, cfg.FirstNamesUrl)
	wg.Wait()

	//using scheduler
	scheduleUpdateLadiesJob(cfg.LadiesCfg, loc, cfg.SocksProxyAddress)

	go manualOrdersJob()
	//everything ready, start sending orders
	go sendOrdersJob(&cfg.OrdersRoutineCfg, loc, cfg.SocksProxyAddress)

	go admin.StartControlPanelServer(&cfg.OrdersRoutineCfg)
}

func manualOrdersJob() {
	log.Info("Starting manual orders RabbitMq listener")
	for {
		func() {
			defer util.RecoverAndLogError("RabbitMq")
			order := rabbitmq.ConsumeManualOrder()
			target.QueueManualOrder(order)
			log.Info("Manual order is queued and will be executed by Orders routine")
		}()
	}
}

func scheduleUpdateLadiesJob(cfg config.LadiesConfig, loc *time.Location, socksProxy string) {
	startTime := time.Now().Add(cfg.UpdateLadiesStartDelay)
	s := gocron.NewScheduler(loc)
	_, err := s.Every(cfg.UpdateLadiesInterval).StartAt(startTime).Do(func() {
		websites.UpdateLadies(cfg.LadiesBaseUrl, cfg.LadiesUrls, socksProxy)
	})

	if err != nil {
		log.WithError(err).Error("Failed to start UpdateLadies job")
	} else {
		log.Info("Starting update ladies routine")
		s.StartAsync()
	}
}

func sendOrdersJob(cfg *config.OrdersRoutineConfig, loc *time.Location, socksProxy string) {
	log.Info("Starting send orders routine")

	if cfg.StartDelay {
		log.Info("Configured to delay the first order")
		sleepDuration := time.Duration(float64(cfg.SendOrdersMaxInterval) * rand.Float64())
		log.Infof("sendOrdersJob: scheduling order in %.2f minutes", float64(sleepDuration)/float64(time.Minute))
		time.Sleep(sleepDuration)
	} else {
		log.Warn("Sending initial order without delay!")
	}

	for {
		log.Info("sendOrdersJob: Order flow triggered")
		sleepAtNight(cfg, loc)

		readyToGo := dao.Dao.ValidateDataIntegrity()
		ordersEnabled := cfg.SendOrdersEnabled

		if readyToGo && ordersEnabled {
			go target.OrderItem(cfg.OrdersCfg, socksProxy)
		} else {
			if !readyToGo {
				log.Warn("Cannot send orders due to empty database tables, please check DB!")
			}
			if !ordersEnabled {
				log.Info("SendOrdersEnabled = false, not sending anything")
			}
		}
		sleepDuration := time.Duration(float64(cfg.SendOrdersMaxInterval) * rand.Float64())
		log.Infof("sendOrdersJob: scheduling next order in %.2f minutes", float64(sleepDuration)/float64(time.Minute))
		time.Sleep(sleepDuration)
	}
}

func sleepAtNight(cfg *config.OrdersRoutineConfig, loc *time.Location) {
	year, month, day := time.Now().In(loc).Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, loc)

	currentTime := time.Now().In(loc)
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
	log.Info("sendOrdersJob: Beyond work hours, sleeping until " +
		time.Now().Add(sleepDuration).Format("2006-01-02 15:04:05"))
	time.Sleep(sleepDuration)
}
