package operator

import (
	"time"

	"github.com/dovics/borja/device/light_sensor"
	"github.com/dovics/borja/util/log"
)

type LightOperator struct {
	sensor *light_sensor.LightSensor
	ch     chan struct{}
	isOpen bool
}

func NewLightOperator(sensor *light_sensor.LightSensor) *LightOperator {
	operator := &LightOperator{
		sensor: sensor,
		ch:     make(chan struct{}),
	}

	if err := sensor.SendAsciiModeChange(); err != nil {
		log.Error(err)
	}

	return operator
}

func (t *LightOperator) Chan() <-chan struct{} {
	return t.ch
}

func (t *LightOperator) QueryLight() (value interface{}, err error) {
	retry := 0
	for {
		value, err = t.sensor.Read()
		if err != nil {
			retry++
			time.Sleep(time.Second)
			if retry >= 20 {
				return
			}

			if retry%10 == 0 {
				if err = t.sensor.Reconnect(); err != nil {
					log.Error("reconnect fail, error: ", err)
				} else {
					log.Info("reconnect success")
				}
			}

			if retry%5 == 0 {
				if err := t.sensor.SendAsciiModeChange(); err != nil {
					log.Error("send mode change fail, error: ", err)
				} else {
					log.Info("send mode change success")
				}
			}
		} else {
			return
		}
	}
}

func (t *LightOperator) IsOpen() bool {
	return t.isOpen
}

func (t *LightOperator) Run() {
	for {
		value, err := t.QueryLight()
		if err != nil {
			log.Error("query light error: ", err)
			continue
		}

		if value.(int) < 1000 && !t.isOpen || value.(int) > 1000 && t.isOpen {
			t.isOpen = !t.isOpen
			select {
			case t.ch <- struct{}{}:
			default:
			}
		}
	}
}
