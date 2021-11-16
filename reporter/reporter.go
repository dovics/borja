package reporter

import (
	"bytes"
	"encoding/json"
	"errors"

	"net/http"

	"github.com/dovics/borja/cache"
	"github.com/dovics/borja/util/log"
)

type Reporter interface {
	Run()
	Register(name string, operate Operate)
	SetTrigger(Trigger) error
}

type CachedReporter interface {
	Reporter
	SetCache(cache.Cache) error
}

type Operate func() (interface{}, error)

type Trigger interface {
	Chan() <-chan struct{}
}

func New(url string) CachedReporter {
	return &defaultReporter{
		target:   url,
		client:   *http.DefaultClient,
		operates: make(map[string]Operate),
	}
}

type defaultReporter struct {
	operates map[string]Operate
	trigger  Trigger
	cache    cache.Cache
	running  bool

	target string
	client http.Client
}

func (r *defaultReporter) Run() {
	r.running = true
	for range r.trigger.Chan() {
		data := make(map[string]interface{})
		var err error
		for k, operate := range r.operates {
			data[k], err = operate()
			if err != nil {
				data[k] = err
			}

			if r.cache != nil {
				err := r.cache.Add(k, data[k])
				if err != nil {
					log.Info("cache save error, key: %s, value: %v, error: %s", k, data[k], err)
				}
			}
		}

		buf, err := json.Marshal(data)
		if err != nil {
			log.Info(err)
			continue
		}

		if r.target != "" {
			resp, err := r.client.Post(r.target, "json", bytes.NewReader(buf))
			if err != nil {
				log.Info(err)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				log.Info(errors.New("resp status isn't OK"))
				continue
			}
		}

	}
}

func (r *defaultReporter) SetTrigger(trigger Trigger) error {
	if r.running {
		return errors.New("reporter is running")
	}
	r.trigger = trigger
	return nil
}

func (r *defaultReporter) SetCache(c cache.Cache) error {
	if r.running {
		return errors.New("reporter is running")
	}

	r.cache = c
	return nil
}

func (r *defaultReporter) Register(name string, operate Operate) {
	r.operates[name] = operate
}
