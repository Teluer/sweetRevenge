package main

import (
	"github.com/magiconair/properties"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"os"
	"sweetRevenge/config"
	"sweetRevenge/db/dao"
	"time"
)

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

	dao.AutoMigrateAll()

	p := properties.MustLoadFile("config.properties", properties.UTF8)
	var cfg config.Config
	if err := p.Decode(&cfg); err != nil {
		log.WithError(err).Fatal("Failed to parse configs")
	}

	programLogic(cfg)

	//wait indefinitely
	select {}
}
