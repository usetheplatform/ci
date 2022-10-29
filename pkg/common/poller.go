package common

import (
	"time"
)

type Poller struct {
	ticker *time.Ticker
	quit   chan bool
}

func NewPoller(callback func(), timeout time.Duration) *Poller {
	poller := &Poller{
		ticker: time.NewTicker(timeout),
	}

	go poller.run(callback)
	return poller
}

func (p *Poller) run(callback func()) {
	for {
		select {
		case <-p.quit:
			p.ticker.Stop()
			return
		case <-p.ticker.C:
			callback()
		}
	}
}

func (p *Poller) Stop() {
	p.quit <- true
}
