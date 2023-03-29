package app

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"time"
)

type (
	// Repo interface for repo data repository.
	Repo interface {
		// SaveOrder adds new order to the repository.
		// Errors: unknown.
		SaveOrder(context.Context, Order) (*Order, error)
		// UpdateOrder order in repo.
		// Errors: unknown.
		UpdateOrder(context.Context, Order) error
		// GetOrder order by id.
		// Errors: ErrNotFound, unknown.
		GetOrder(context.Context, uuid.UUID) (*Order, error)
		// ListOrders returns orders by params.
		// Errors: unknown.
		ListOrders(context.Context, OrderParams) ([]Order, int, error)
	}

	// TaskRepo interface for saving orders.
	TaskRepo interface {
		// SaveTask adds new task to repository.
		// Errors: unknown.
		SaveTask(context.Context, Task) (uuid.UUID, error)
		// FinishTask set column Task.FinishedAt task.
		// Errors: ErrNotFound, unknown.
		FinishTask(context.Context, uuid.UUID) error
		// ListActualTask returns list task by limit and ordered by created_at (ask).
		// Return tasks without Task.FinishedAt.
		// Errors: unknown.
		ListActualTask(context.Context, int) ([]Task, error)
	}

	// Queue sends events to queue.
	Queue interface {
		// AddOrder sends event 'EventAddOrder' to queue.
		// Errors: unknown.
		AddOrder(context.Context, EventAddOrder) error
		// UpdateOrder sends event 'EventUpdateOrder' to queue.
		// Errors: unknown.
		UpdateOrder(ctx context.Context, eventUpdate EventUpdateOrder) error
	}

	Cron interface {
		// Fetch get jobs.
		// Errors: unknown.
		Fetch() <-chan time.Time
	}
)
