package light_sensor

import (
	"testing"
	"time"

	"github.com/tarm/serial"
)

func TestLightSensorBySerial(t *testing.T) {
	c := &serial.Config{Name: "/dev/serial0", Baud: 9600, ReadTimeout: time.Second * 5}
	
	sensor, err := ConnectBySerial(c)
	if err != nil {
		t.Fatal(err)
	}
	for {
		data, err := sensor.Read()
		if err != nil {
			t.Error(err)
		}

		t.Log(data)
	}

}

func TestLightSensorModeChange(t *testing.T) {
	c := &serial.Config{Name: "/dev/serial0", Baud: 9600, ReadTimeout: time.Second * 5}
	
	sensor, err := ConnectBySerial(c)
	if err != nil {
		t.Fatal(err)
	}

	if err := sensor.SendAsciiModeChange();err != nil {
		t.Fatal(err)
	}
}