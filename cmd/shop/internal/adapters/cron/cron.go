package cron

import (
	"context"
	"github.com/robfig/cron/v3"
	"time"
)

type (
	Config struct {
		TimeFetch string
	}

	Cron struct {
		timeouts     timeouts
		scheduler    *cron.Cron
		chEventFetch chan time.Time
	}

	timeouts struct {
		Fetch string
		Send  string
	}
)

func New(cfg Config) *Cron {
	return &Cron{
		timeouts: timeouts{
			Fetch: cfg.TimeFetch,
		},
		scheduler:    cron.New(),
		chEventFetch: make(chan time.Time),
	}
}

// TimeFetch implements app.Cron.
func (s *Cron) TimeFetch() <-chan time.Time {
	return s.chEventFetch
}

func (s *Cron) Process(ctx context.Context) (err error) {
	_, err = s.scheduler.
		AddFunc(s.timeouts.Fetch, func() {
			// just for an example
			t := time.Now()
			s.chEventFetch <- t
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
