package app

import (
	"context"
	"fmt"
)

func (a *App) Process(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case limit := <-a.cron.Fetch():
			fmt.Println(limit)
			// logic for fetch tasks and send it in queue
		}
	}
}

func (a *App) Events() ([]Task, error) {

}
