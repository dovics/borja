package light_sensor

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"

	"github.com/dovics/borja/util/log"
	"github.com/tarm/serial"
)

var ErrWrongReply = errors.New("wrong reply")

var (
	assicReg       *regexp.Regexp
	okReg          *regexp.Regexp
	okAndLightReg  *regexp.Regexp
	initSuccessReg *regexp.Regexp
)

func init() {
	assicReg = regexp.MustCompile(`\s*Ligth:(\d+)\s*lx\s*`)
	okReg = regexp.MustCompile(`\s*OK\s*`)
	okAndLightReg = regexp.MustCompile(`\s*OK\s*Light:(\d+)\s*lx`)
	initSuccessReg = regexp.MustCompile(`\s*INIT\s*SUCCES\s*`)
}

type LightSensor struct {
	config *serial.Config
	port   *serial.Port
	buf    *bufio.Reader
}

func ConnectBySerial(config *serial.Config) (*LightSensor, error) {
	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	return &LightSensor{
		config: config,
		port:   port,
		buf:    bufio.NewReader(port),
	}, nil
}

func (s *LightSensor) Reconnect() error {
	port, err := serial.OpenPort(s.config)
	if err != nil {
		return err
	}

	s.port = port
	s.buf = bufio.NewReader(port)
	return nil
}

func (s *LightSensor) Read() (int, error) {
	line, _, err := s.buf.ReadLine()
	if err != nil {
		return 0, err
	}

	result := assicReg.FindSubmatch(line)
	if len(result) > 0 {
		return strconv.Atoi(string(result[1]))
	}

	return 0, errors.New("can't find match")
}

func (s *LightSensor) IsConnected() bool {
	_, err := s.port.Write([]byte("AT"))
	if err != nil {
		return false
	}

	if !s.waitReply(okReg) {
		return false
	}

	return true
}

func (s *LightSensor) SendChipModeChange() error {
	return s.SetFeedBackSpeed(0)
}

func (s *LightSensor) SendAsciiModeChange() error {
	_, err := s.port.Write([]byte("AT+INIT"))
	if err != nil {
		return err
	}

	// if !s.waitReply(initSuccessReg) {
	// 	log.Error(ErrWrongReply)
	// }

	return s.SetFeedBackSpeed(1000)
}

func (s *LightSensor) SetFeedBackSpeed(speed int) error {
	_, err := s.port.Write([]byte("AT+PRATE=" + strconv.Itoa(speed)))
	if err != nil {
		return err
	}

	// if !s.waitReply(okAndLightReg) {
	// 	return ErrWrongReply
	// }

	return nil
}

func (s *LightSensor) waitReply(expect *regexp.Regexp) bool {
	retryTimes := 0

retry:
	line, _, err := s.buf.ReadLine()
	if err != nil {
		log.Error(err)
		if err == io.EOF && retryTimes < 3 {
			retryTimes++
			goto retry
		}
		return false
	}

	if !expect.Match(line) {
		log.Debug("wrong reply: ", string(line))
		return false
	}

	return true
}
