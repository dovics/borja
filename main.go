package main

import (
	"log"
	"time"

	"github.com/tarm/serial"

	"github.com/dovics/borja/cache"
	"github.com/dovics/borja/device/light_sensor"
	"github.com/dovics/borja/exporter"
	"github.com/dovics/borja/operator"
	"github.com/dovics/borja/reporter"
	"github.com/dovics/borja/util/trigger"
)

const (
	boltdbName  = "bolt"
	archivePath = "./data"
)

func main() {
	config := &serial.Config{Name: "COM4", Baud: 9600, ReadTimeout: time.Second * 5}

	sensor, err := light_sensor.ConnectBySerial(config)
	if err != nil {
		log.Fatal(err)
	}
	lightOperator := operator.NewLightOperator(sensor)

	c, err := cache.NewBlotCache(boltdbName)
	if err != nil {
		log.Fatal(err)
	}
	cache.AutoArchiveWrapper(c, archivePath, time.Minute)
	reporter := reporter.New("")

	if err := reporter.SetTrigger(trigger.NewTimeTrigger(time.Second)); err != nil {
		log.Fatal(err)
	}

	if err := reporter.SetCache(c); err != nil {
		log.Fatal(err)
	}

	reporter.Register("light", lightOperator.QueryLight)
	go reporter.Run()

	exporter := exporter.NewFileExporter(":8080")
	exporter.Register("data", func() (interface{}, error) {
		return archivePath, nil
	})

	exporter.Run()
}
