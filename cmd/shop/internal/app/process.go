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
		case t := <-a.cron.Fetch():
			fmt.Println(t)
			// logic for fetch tasks and send it in queue
		}
	}
}
