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

const boltdbName = "bolt"

func main() {
	config := &serial.Config{Name: "COM4", Baud: 9600, ReadTimeout: time.Second * 5}

	sensor, err := light_sensor.ConnectBySerial(config)
	if err != nil {
		log.Fatal(err)
	}
	lightOperator := operator.NewLightOperator(sensor)

	c, err := cache.NewCache(boltdbName)
	if err != nil {
		log.Fatal(err)
	}

	reporter := reporter.New("")

	if err := reporter.SetTrigger(trigger.NewTimeTrigger(time.Second)); err != nil {
		log.Fatal(err)
	}

	if err := reporter.SetCache(c); err != nil {
		log.Fatal(err)
	}

	reporter.Register("light", lightOperator.QueryLight)
	go reporter.Run()

	exporter := exporter.NewExporter()
	exporter.Register("light", lightOperator.QueryLight)
	exporter.Register("data", func() (interface{}, error) {
		data, err := c.GetAfter(time.Now().Add(-time.Minute * 15))
		if err != nil {
			return nil, err
		}

		if err := c.ClearBefore(time.Now()); err != nil {
			return nil, err
		}

		return data, nil
	})

	exporter.Run()
}
