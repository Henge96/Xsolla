package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
	"xsolla/internal/dom"
)

func (a *App) Process(ctx context.Context) error {
	for {
		var (
			err   error
			tasks = make([]Task, 0)
		)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case limit := <-a.cron.Fetch():
			tasks, err = a.repo.ListActualTask(ctx, int(limit))
			if err != nil {
				log.Println("couldn't get tasks", err)

				continue
			}

			if len(tasks) == 0 {
				continue
			}

			for i := range tasks {
				switch tasks[i].TaskKind {
				case dom.TaskKindEventAdd:
					err = a.handleTaskKindEventAdd(ctx, tasks[i])
				case dom.TaskKindEventUpdate:
					err = a.handleTaskKindEventUpdate(ctx, tasks[i])
				default:
					log.Println("unknown task kind", tasks[i])

					continue
				}
				if err != nil {
					log.Println("couldn`t send event in queue", tasks[i])

					continue
				}
			}
		case msg := <-a.queue.UpdateOrderStatus():
			err = a.handleUpdateOrderStatus(ctx, msg)
			if err != nil {
				log.Println("couldn't handle event", err)

				continue
			}
		}
	}
}

func (a *App) collectingTasks(ctx context.Context, wg *sync.WaitGroup, out chan Task, taskLimit int) {
	defer wg.Done()

	const (
		taskTickerTimeout = time.Second / 10
	)

	var (
		ticker = time.NewTicker(taskTickerTimeout)
	)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return
	case <-ticker.C:
		tasks, err := a.repo.ListActualTask(ctx, taskLimit)
		if err != nil {
			log.Println("couldn't get tasks", err)

			return
		}

		for i := range tasks {
			select {
			case <-ctx.Done():
				return
			case out <- tasks[i]:
			}
		}
	}

}

func (a *App) handleTaskKindEventAdd(ctx context.Context, task Task) error {
	err := a.queue.AddOrder(ctx, EventAddOrder{
		TaskID:    task.ID,
		Status:    task.Order.Status,
		CreatedAt: task.Order.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("a.queue.AddComment: %w", err)
	}

	// TODO validation mechanism for check that task will finished
	err = a.repo.FinishTask(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("a.repo.FinishTask: %w", err)
	}

	return nil
}

func (a *App) handleTaskKindEventUpdate(ctx context.Context, task Task) error {
	err := a.queue.UpdateOrder(ctx, EventUpdateOrder{
		TaskID:    task.ID,
		ID:        task.Order.ID,
		Status:    task.Order.Status,
		CreatedAt: task.Order.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("a.queue.AddComment: %w", err)
	}

	// TODO validation mechanism for check that task will finished
	err = a.repo.FinishTask(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("a.repo.FinishTask: %w", err)
	}

	return nil
}

func (a *App) handleUpdateOrderStatus(ctx context.Context, event dom.Event[EventUpdateOrderStatusFromQueue]) error {
	order := Order{
		ID:        event.Body().SourceID,
		Status:    event.Body().Status,
	}

	_, err := a.repo.UpdateOrder(ctx, order)
	switch {
	case errors.Is(err, ErrDuplicate):
	// We must acknowledge this message.
	case err != nil:
		event.Nack(ctx)

		return fmt.Errorf("a.repo.SavePost: %w", err)
	}

	event.Ack(ctx)

	return nil
}
