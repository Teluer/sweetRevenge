package main

import (
	"github.com/magiconair/properties"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"os"
	"sweetRevenge/config"
	"time"
)

func init() {
	//log.SetReportCaller(true)
	log.Info("Program Startup")

	file, err := os.Create("sweetRevenge.log")
	if err != nil {
		log.Fatal("failed to create log file, what's the point now...")
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))

	rand.Seed(time.Now().UnixMilli())
}

func main() {
	p := properties.MustLoadFile("config.properties", properties.UTF8)
	var cfg config.Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	//programLogic(cfg)
	//test.TestAnonSending()
	//test.SendTestRequest()
	//websites.UpdateLadies()
	//target.ExecuteManualOrder()

	//wait indefinitely
	select {}
}
