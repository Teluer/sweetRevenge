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
	"sweetRevenge/src/util"
	"time"
)

// TODO: add proper documentation for some stuff
func main() {
	log.Info("Program startup")
	rand.Seed(time.Now().UnixMilli())

	//load configs
	p := properties.MustLoadFile("config.properties", properties.UTF8)
	var cfg config.Config
	if err := p.Decode(&cfg); err != nil {
		log.WithError(err).Fatal("Failed to parse configs")
	}

	//create a log file
	file, err := os.Create("sweetRevenge.log")
	if err != nil {
		log.WithError(err).Fatal("failed to even create log file, what's the point now...")
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	loc, _ := time.LoadLocation(cfg.TimeZone)
	log.SetFormatter(util.LogFormatter{Formatter: &log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"}, Loc: loc})

	dao.Dao.OpenDatabaseConnection(cfg.DatabaseDsn)
	dao.Dao.AutoMigrateAll()
	rabbitmq.InitializeRabbitMq(cfg.Rabbit)
	programLogic(&cfg, loc)

	//wait indefinitely
	select {}
}
