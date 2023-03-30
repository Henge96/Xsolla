package app

import (
	"context"
	"fmt"
)

func (a *App) Process(ctx context.Context) (err error) {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case limit := <-a.cron.Fetch():
			fmt.Println(limit)
			_, err = a.Events(ctx, limit)
			// logic for fetch tasks and send it in queue
		}
	}
}

func (a *App) Events(ctx context.Context, limit uint8) ([]Task, error) {
	_, _ = a.repo.ListActualTask(ctx, int(limit))
	return nil, nil
}
