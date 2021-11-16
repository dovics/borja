package trigger

import "time"

type TimeTrigger struct {
	ticker *time.Ticker
	ch     chan struct{}
	stop   chan struct{}
}

func NewTimeTrigger(d time.Duration) *TimeTrigger {
	ticker := time.NewTicker(d)
	trigger := &TimeTrigger{
		ticker: ticker,
		ch:     make(chan struct{}),
		stop:   make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-ticker.C:
				trigger.ch <- struct{}{}
			case <-trigger.stop:
				return
			}
		}
	}()

	return trigger
}

func (t *TimeTrigger) Chan() <-chan struct{} {
	return t.ch
}

func (t *TimeTrigger) Stop() {
	close(t.stop)
	t.ticker.Stop()
}
