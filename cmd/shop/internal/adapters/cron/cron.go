package cron

import (
	"context"
	"github.com/robfig/cron/v3"
)

type (
	Config struct {
		TimeFetch string
		Limit     uint8
	}

	Cron struct {
		fetch     fetch
		scheduler *cron.Cron
	}

	fetch struct {
		Timeout      string
		Limit        uint8
		chEventFetch chan uint8
	}
)

func New(cfg Config) *Cron {
	return &Cron{
		fetch: fetch{
			Timeout:      cfg.TimeFetch,
			Limit:        cfg.Limit,
			chEventFetch: make(chan uint8),
		},
		scheduler: cron.New(),
	}
}

// Fetch implements app.Cron.
func (s *Cron) Fetch() <-chan uint8 {
	return s.fetch.chEventFetch
}

func (s *Cron) Process(ctx context.Context) (err error) {
	_, err = s.scheduler.
		AddFunc(s.fetch.Timeout, func() {
			s.fetch.chEventFetch <- s.fetch.Limit
		})
	if err != nil {
		return
	}

	go func() {
		select {
		case <-ctx.Done():
			s.scheduler.Stop()
			return
		}
	}()
	s.scheduler.Start()
	return nil
}
