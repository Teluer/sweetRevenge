package main

import (
	"github.com/magiconair/properties"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"os"
	"sweetRevenge/src/config"
	"sweetRevenge/src/db/dao"
	"sweetRevenge/src/rabbitmq"
	"time"
)

//TODO: develop admin panel to update several configs
//TODO: pass manual orders through an MQ
//TODO: add unit testing

func init() {
	file, err := os.Create("sweetRevenge.log")
	if err != nil {
		log.Fatal("failed to even create log file, what's the point now...")
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))

	rand.Seed(time.Now().UnixMilli())
}

func main() {
	log.Info("Program Startup")

	p := properties.MustLoadFile("config.properties", properties.UTF8)
	var cfg config.Config
	if err := p.Decode(&cfg); err != nil {
		log.WithError(err).Fatal("Failed to parse configs")
	}
	
	//TODO: not smart to keep one connection for the entire lifecycle
	dao.Dao.OpenDatabaseConnection()
	dao.Dao.AutoMigrateAll()

	rabbitmq.InitializeRabbitMq(cfg.Rabbit)

	programLogic(&cfg)

	//wait indefinitely
	select {}
}
