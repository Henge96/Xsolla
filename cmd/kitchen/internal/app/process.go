package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"log"
	"math/rand"
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
			err = a.handleUpdateOrder(ctx, msg)
			if err != nil {
				log.Println("couldn't handle event", err)

				continue
			}
		case msg := <-a.queue.AddOrder():
			err = a.handleNewOrder(ctx, msg)
			if err != nil {
				log.Println("couldn't handle event", err)

				continue
			}
		}
	}
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

func (a *App) handleNewOrder(ctx context.Context, event dom.Event[EventAddOrderFromQueue]) error {
	rand.Seed(time.Now().Unix())

	err := a.repo.Tx(ctx, func(repo Repo) error {
		order, err := a.repo.SaveOrder(ctx, event.Body().Order)
		switch {
		case errors.Is(err, ErrDuplicate):
			// We must acknowledge this message.
			event.Ack(ctx)

			return nil
		case err != nil:
			return fmt.Errorf("a.repo.SaveOrder: %w", err)
		}

		c := Cooking{
			Status: CookingStatusNew,
			// in production project get from repo
			Chef:    uuid.Must(uuid.NewV4()),
			Table:   rand.Uint32(),
			OrderID: order.ID,
		}

		_, err = a.repo.SaveCooking(ctx, c)
		if err != nil {
			return fmt.Errorf("a.repo.SaveCooking: %w", err)
		}

		return nil
	})
	if err != nil {
		event.Nack(ctx)

		return err
	}

	event.Ack(ctx)

	return nil
}

func (a *App) handleUpdateOrder(ctx context.Context, event dom.Event[EventUpdateOrderStatusFromQueue]) error {
	// todo add logic if order was canceled in cooking process
	if event.Body().Status == dom.OrderStatusConfirmed {
		order, err := a.repo.GetOrder(ctx, event.Body().SourceID)
		switch {
		case errors.Is(err, ErrNotFound):
			// todo think how to fix or notify about this
			event.Ack(ctx)

			return nil
		case err != nil:
			event.Nack(ctx)

			return fmt.Errorf("a.repo.GetPost: %w", err)
		}

		order.Status = event.Body().Status

		err = a.repo.Tx(ctx, func(repo Repo) error {
			_, err = a.repo.UpdateOrder(ctx, *order)
			if err != nil {
				return fmt.Errorf("a.repo.UpdateOrder: %w", err)
			}

			_, err = a.repo.UpdateCookingStatusByOrderID(ctx, order.ID, CookingStatusNeedToStart)
			if err != nil {
				return fmt.Errorf("a.repo.UpdateCookingStatusByOrderID: %w", err)
			}

			return nil
		})
		if err != nil {
			event.Nack(ctx)

			return fmt.Errorf("a.repo.GetPost: %w", err)
		}
	}
	event.Ack(ctx)

	return nil
}
